package clirpc

import (
	"fmt"
	"log"
  "net"
	"net/rpc"
	"testing"
	"time"
)

const (
	srvAddr = "127.0.0.1:58085"
)

type Listener struct {
	Sleep time.Duration
}

func (l *Listener) GetUser(line []byte, ack *bool) (err error) {
	fmt.Println(string(line))
	return
}

func startServer(stop,started chan bool) {
	fmt.Println("Started server")

  addr, err := net.ResolveTCPAddr("tcp", srvAddr)
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	listener := new(Listener)
	rpc.Register(listener)
	fmt.Println("server")
	go rpc.Accept(inbound)
  started<-true
  <-stop
}

func TestServer(t *testing.T) {
	var stopSrv chan bool
	stopSrv = make(chan bool)
  srvStarted := make(chan bool)
	go startServer(stopSrv, srvStarted)
  <-srvStarted

	client, err := rpc.Dial("tcp", srvAddr)
	if err != nil {
		log.Fatal(err)
	}

	var reply bool
	var line []byte
	line = []byte("200024")
	err = client.Call("Listener.GetUser", line, &reply)
	if err != nil {
		log.Fatal(err)
	}

	stopSrv <- true
	fmt.Println("testing done")
}
