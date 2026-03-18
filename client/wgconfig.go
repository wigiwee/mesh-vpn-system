package main

import (
	"client/config"
	"client/models"
	"fmt"
	"os"
	"strings"
)

func writeWGConfig(privateKey string, selfIP string, peers []models.Peer) error {

	var sb strings.Builder
	sb.WriteString("[interface]\n")
	sb.WriteString(fmt.Sprintf("PrivateKey = %s\n", privateKey))
	sb.WriteString(fmt.Sprintf("Address = %s/32\n", selfIP))
	sb.WriteString("ListenPort = 51820\n\n")

	for _, peer := range peers {

		sb.WriteString("[Peer]\n")
		sb.WriteString(fmt.Sprintf("PublicKey = %s\n", peer.PublicKey))
		sb.WriteString(fmt.Sprintf("AllowedIPs = %s/32\n", peer.IPAddress))
		sb.WriteString(fmt.Sprintf("Endpoint = %s\n\n", peer.Endpoint))
	}

	return os.WriteFile(config.WG_CONFIG_FILE_LOCATION, []byte(sb.String()), 0600)
}
