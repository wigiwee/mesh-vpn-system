package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func main() {

	privateKey, publicKey := generateKeys()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	template, _ := os.ReadFile("../scripts/wgtemplate.conf")
	config := string(template)

	config = strings.ReplaceAll(config, "{{PRIVATE_KEY}}", privateKey)
	config = strings.ReplaceAll(config, "{{ADDRESS}}", "10.0.0.1/24")
	config = strings.ReplaceAll(config, "{{PEER_PUBLIC_KEY}}", "<PUT_PEER_PUBLIC_KEY_HERE>")
	config = strings.ReplaceAll(config, "{{PEER_ALLOWED_IPS}}", "10.0.0.2/32")
	config = strings.ReplaceAll(config, "{{PEER_ENDPOINT}}", "<PUT_PEER_IP>:51820")

	err := os.WriteFile("/etc/wireguard/wg0.conf", []byte(config), 0600)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Saved config to wg0.conf ✅")

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
