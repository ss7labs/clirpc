package clirpc

import (
	"context"
	"net"
	"sync"
	"testing"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

type testRadiusServer struct {
	server   *radius.PacketServer
	secret   []byte
	listener net.PacketConn
	wg       sync.WaitGroup
}

func newTestRadiusServer(t *testing.T, addr string, secret string) *testRadiusServer {
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		t.Fatalf("failed to create test server: %v", err)
	}

	server := &testRadiusServer{
		secret:   []byte(secret),
		listener: pc,
	}

	handler := radius.HandlerFunc(func(w radius.ResponseWriter, r *radius.Request) {
		username := rfc2865.UserName_GetString(r.Packet)
		if username != "0:0:130" {
			w.Write(r.Response(radius.CodeDisconnectNAK))
			return
		}

		w.Write(r.Response(radius.CodeDisconnectACK))
	})

	server.server = &radius.PacketServer{
		Handler:      handler,
		SecretSource: radius.StaticSecretSource(server.secret),
	}

	return server
}

func (s *testRadiusServer) start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.server.Serve(s.listener); err != radius.ErrServerShutdown {
			panic(err)
		}
	}()
}

func (s *testRadiusServer) stop() {
	s.server.Shutdown(context.Background())
	s.wg.Wait()
}

func TestSendRadiusDisconnectWithServer(t *testing.T) {
	// Start test server
	server := newTestRadiusServer(t, "127.0.0.1:0", "testing123")
	server.start()
	defer server.stop()

	// Get the actual address the server is listening on
	addr := server.listener.LocalAddr().(*net.UDPAddr)

	// Test cases
	tests := []struct {
		name     string
		config   RadiusConfig
		wantErr  bool
		username string
	}{
		{
			name: "valid disconnect request",
			config: RadiusConfig{
				Host:     "127.0.0.1",
				Port:     addr.Port,
				Secret:   "testing123",
				Username: "0:0:130",
			},
			wantErr: false,
		},
		{
			name: "invalid username",
			config: RadiusConfig{
				Host:     "127.0.0.1",
				Port:     addr.Port,
				Secret:   "testing123",
				Username: "invalid",
			},
			wantErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendRadiusDisconnect(context.Background(), tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendRadiusDisconnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
