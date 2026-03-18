package config

import (
	"client/models"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func WriteWGConfig() error {

	var sb strings.Builder
	sb.WriteString("[interface]\n")
	sb.WriteString(fmt.Sprintf("PrivateKey = %s\n", ConfigObj.PrivateKey))
	sb.WriteString(fmt.Sprintf("Address = %s/32\n", ConfigObj.NodeIPAddr))
	sb.WriteString("ListenPort = 51820\n\n")

	// for _, peer := range peers {

	// 	sb.WriteString("[Peer]\n")
	// 	sb.WriteString(fmt.Sprintf("PublicKey = %s\n", peer.PublicKey))
	// 	sb.WriteString(fmt.Sprintf("AllowedIPs = %s/32\n", peer.IPAddress))
	// 	sb.WriteString(fmt.Sprintf("Endpoint = %s\n\n", peer.Endpoint))
	// }

	return os.WriteFile(WG_CONFIG_FILE_LOCATION, []byte(sb.String()), 0600)
}

func AddPeer(p models.Peer) error {
	exec.Command("wg", "set", "wg0",
		"peer", p.PublicKey,
		"allowed-ips", p.IPAddress+"/32",
		"endpoint", p.Endpoint,
	).Run()
	return nil
}

func RemovePeer(p models.Peer) {
	exec.Command("wg", "set", "wg0",
		"peer", p.PublicKey,
		"remove",
	).Run()
}
