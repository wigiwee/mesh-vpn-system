package main

import (
	"client/config"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/pion/stun"
)

func GetPublicEndpoint() (string, error) {
	conn, err := net.Dial("udp", config.STUN_SERVERS[1])
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client, err := stun.NewClient(conn)
	if err != nil {
		return "", err
	}
	defer client.Close()
	done := make(chan struct{})
	var publicAddr string

	err = client.Do(stun.MustBuild(stun.TransactionID, stun.BindingRequest), func(e stun.Event) {
		if e.Error != nil {
			log.Println("stun error " + err.Error())
			close(done)
			return
		}

		var xorAddr stun.XORMappedAddress
		if getErr := xorAddr.GetFrom(e.Message); getErr != nil {
			log.Println("Parse Error " + getErr.Error())
			close(done)
			return
		}
		publicAddr = xorAddr.IP.String() + ":" + strconv.Itoa(xorAddr.Port)
		close(done)
	})
	if err != nil {
		return "", err
	}
	select {
	case <-done:
		return publicAddr, nil
	case <-time.After(3 * time.Second):
		return "", fmt.Errorf("STUN timeout")
	}

}
