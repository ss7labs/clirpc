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

var testData [2]RawSession

func initTestData() {
	testData[0] = RawSession{"25330868", "737155@gts", "7C:8B:CA:E2:DB:BE", "1", "3065", "1031", "1669976882-25330868", "172.17.5.62", "1492", "8M;10M", "2M;15M", "62299192", "56293332", "22641605997", "40947477724", "10 day(s), 5 hour(s), 46 min(s), 16 sec(s)"}
	testData[1] = RawSession{"37905591", "430147@gts", "84:A9:C4:F9:56:A3", "1", "3131", "1061", "1670854858-37905591", "172.17.40.29", "1480", "-", "-", "0", "0", "0", "0", "0 day(s), 1 hour(s), 53 min(s), 20 sec(s)"}
}

type Listener struct {
	Sleep time.Duration
}

func (l *Listener) GetUser(line []byte, ack *bool) (err error) {
	fmt.Println(string(line))
	return
}

func startServer(stop, started chan bool) {

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
	fmt.Println("Started server")
	go rpc.Accept(inbound)
	started <- true
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
