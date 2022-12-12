package clirpc

import (
	"fmt"
	"log"
	"net/rpc"
	"testing"
	"time"
)

const (
	srvAddr = "127.0.0.1:58085"
)

func startServer(stop chan bool) {
	fmt.Println("Started server")
	<-stop
}

func TestServer(t *testing.T) {
	var stopSrv chan bool
	stopSrv = make(chan bool)
	go startServer(stopSrv)

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

	time.Sleep(1 * time.Second)
	stopSrv <- true
	fmt.Println("testing done")
}
