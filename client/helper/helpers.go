package helper

import (
	"client/api"
	"client/config"
	"math/rand/v2"
	"strconv"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func GenerateKeys() (private, public string) {
	privateKey, _ := wgtypes.GeneratePrivateKey()
	return privateKey.String(), privateKey.PublicKey().String()
}

func GenerateRandomIPAddr() string {

	ipAddr := "100."
	ipAddr = ipAddr + strconv.Itoa(rand.IntN(255)) + "."
	ipAddr = ipAddr + strconv.Itoa(rand.IntN(255)) + "."
	ipAddr = ipAddr + strconv.Itoa(rand.IntN(255))
	return ipAddr
}

func GetIpAddr() (string, error) {

	ipAddr := GenerateRandomIPAddr()
	nodes, err := api.GetPeers(config.USER_ID, config.USER_ID)
	if err != nil {
		return "", err
	}
	for {
		broke := false
		for _, node := range nodes {
			if node.IPAddress == ipAddr {
				ipAddr = GenerateRandomIPAddr()
				broke = true
				break
			}
		}
		if !broke {
			break
		}
	}
	return ipAddr, nil
}
