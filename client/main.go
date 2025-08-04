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
	fmt.Println(privateKey)
	fmt.Println(publicKey)

	selfIf := "10.0.0.1"
	endpoint := "192.168.0.5:51820"

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
