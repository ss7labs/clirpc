package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
  "fmt"
  "clirpc"
  "net"
  "net/rpc"
)

func main() {
	err := godotenv.Load("tr.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
  port := os.Getenv("PORT")
  addr := os.Getenv("ADDR")

  bind := addr+":"+port
  bindaddr, err := net.ResolveTCPAddr("tcp", bind)
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", bindaddr)
	if err != nil {
		log.Fatal(err)
	}

	listener := new(clirpc.Listener)
	rpc.Register(listener)
	fmt.Println("Started server for ", bind)
	rpc.Accept(inbound)
}

