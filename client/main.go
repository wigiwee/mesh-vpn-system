package main

import (
	"client/api"
	"client/config"
	"client/helper"
	"client/models"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {

	privateKey, publicKey := helper.GenerateKeys()

	selfIp, err := helper.GetIpAddr()
	if err != nil {
		log.Fatal(err)
	}
	endpoint, err := GetPublicEndpoint()
	if err != nil {
		log.Println("STUN failed: " + err.Error())
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("error fetching hostname")
	}
	nodeId, err := api.RegisterNode(models.RegisterRequest{
		PublicKey: publicKey,
		IPAddress: selfIp,
		Endpoint:  endpoint,
		Device:    hostname,
		UserId:    config.USER_ID,
	})
	if err != nil {
		log.Panic(err)
	}
	config.NODE_ID = nodeId[1 : len(nodeId)-2]
	fmt.Println("nodeID ", config.NODE_ID)
	peers, err := api.GetPeers(config.USER_ID, config.NODE_ID)
	if err != nil {
		log.Panic(err)
	}
	err = writeWGConfig(privateKey, selfIp, peers)
	if err != nil {
		log.Panic(err)
	}

	cmd := exec.Command("sudo", "wg-quick", "up", "wg0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	fmt.Println("Brought up WireGuard interface ✅")
}
