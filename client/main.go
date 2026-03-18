package main

import (
	"client/api"
	"client/config"
	"client/models"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"os/exec"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func main() {

	// TODO: check auth first
	config.ReadConfigFile()
	if config.ConfigObj.UserId == "" {
		//TODO: register/authenticate user
		config.ConfigObj.UserId = "6893814a3b3b86cffb0eaea1"
	}
	if config.ConfigObj.NodeId == "" {
		//TODO: register node
		//fetching endpoint
		endpoint, err := GetPublicEndpoint()
		if err != nil {
			log.Println("STUN failed: " + err.Error())
		}
		config.ConfigObj.Endpoint = endpoint
		//fetching hostname
		hostname, err := os.Hostname()
		if err != nil {
			log.Println("error fetching hostname")
			config.ConfigObj.Hostname = rand.Text()[:6] //assigning random hostname
		} else {
			config.ConfigObj.Hostname = hostname
		}
		//generating public private kyes
		privateKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			log.Panic(err)
		}
		config.ConfigObj.PrivateKey, config.ConfigObj.PublicKey = privateKey.String(), privateKey.PublicKey().String()
		//registering node
		registerNodeRes, err := api.RegisterNode(models.RegisterNodeRequest{
			PublicKey: config.ConfigObj.PublicKey,
			Endpoint:  config.ConfigObj.Endpoint,
			Device:    config.ConfigObj.Hostname,
			UserId:    config.ConfigObj.UserId,
			Hostname:  config.ConfigObj.Hostname,
		})
		fmt.Println(registerNodeRes)
		if err != nil {
			log.Panic(err)
		}
		config.ConfigObj.NodeId = registerNodeRes.NodeId
		config.ConfigObj.NodeIPAddr = registerNodeRes.IPAddress

		//TODO: write the created configObj and write to file
		config.WriteConfigFile()
	}

	peers, err := api.GetPeers(config.ConfigObj.UserId, config.ConfigObj.NodeId)
	if err != nil {
		log.Panic(err)
	}
	err = writeWGConfig(config.ConfigObj.PrivateKey, config.ConfigObj.NodeIPAddr, peers)
	if err != nil {
		log.Panic(err)
	}
	config.WriteConfigFile()
	initialcmd := exec.Command("sudo", "wg-quick", "down", config.WG_CONFIG_FILE_LOCATION)
	initialcmd.Stdout = os.Stdout
	initialcmd.Stderr = os.Stderr
	err = initialcmd.Run()
	if err != nil {
		log.Printf("[ERROR]: %s\n", err)
	}
	cmd := exec.Command("sudo", "wg-quick", "up", config.WG_CONFIG_FILE_LOCATION)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("[ERROR]: %s\n", err)
	}
	fmt.Println("Brought up WireGuard interface ✅")
}
