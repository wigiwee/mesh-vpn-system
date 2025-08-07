package main

import (
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	Addr      *net.UDPAddr
	LastSeen  time.Time
	PublicKey string
}

var (
	clients   = make(map[string]*Client)
	clientsMu sync.Mutex
)

func main() {

	addr, err := net.ResolveUDPAddr("udp", ":3478")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("server started on port 3478")

	buf := make([]byte, 65535)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
	}
}
