package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func main() {

	privateKey, publicKey := generateKeys()

	selfIf := "10.0.0.1"
	endpoint, err := GetPublicEndpoint()
	if err != nil {
		log.Println("STUN failed: " + err.Error())
	}

	nodeId, err := registerNode(RegisterRequest{
		PublicKey: publicKey,
		IPAddress: selfIf,
		Endpoint:  endpoint,
		Device:    "deviceName",
		UserId:    USER_ID,
	})
	if err != nil {
		log.Panic(err)
	}
	NODE_ID = nodeId[1 : len(nodeId)-2]
	fmt.Println("nodeID ", NODE_ID)
	peers, err := getPeers(USER_ID, NODE_ID)
	if err != nil {
		log.Panic(err)
	}
	err = writeWGConfig(privateKey, selfIf, peers)
	if err != nil {
		log.Panic(err)
	}

	cmd := exec.Command("sudo", "wg-quick", "up", "wg0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	fmt.Println("Brought up WireGuard interface ✅")
}

func generateKeys() (private, public string) {
	privateKey, _ := wgtypes.GeneratePrivateKey()
	return privateKey.String(), privateKey.PublicKey().String()
}
