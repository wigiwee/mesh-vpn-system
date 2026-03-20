package config

import (
	"client/models"
	"encoding/json"
	"log"
	"os"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const (
	WG_CONFIG_FILE_LOCATION  = CONFIG_DIR + "/" + INTERFACE_NAME + ".conf"
	SERVER_URL               = "http://localhost:4000"
	CONFIG_DIR               = "local_files"
	APP_CONFIG_FILE_LOCATION = CONFIG_DIR + "/" + "app.json"
	INTERFACE_NAME           = "wg0"
)

var (
	STUN_SERVERS []string = []string{
		"stun.l.google.com:19302",
		"stun1.l.google.com:19302",
		"stun2.l.google.com:19302",
		"stun3.l.google.com:19302",
		"stun4.l.google.com:19302",
	}
	ConfigObj Config
	Peers     = make(map[string]models.Peer)
)

type Config struct {
	PublicKey  string
	PrivateKey string
	UserId     string
	NodeId     string
	NodeIPAddr string
	Endpoint   string
	Hostname   string
}

func WriteConfigFile() error {
	configFile, err := os.OpenFile(APP_CONFIG_FILE_LOCATION, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(ConfigObj)
	if err != nil {
		return err
	}
	return nil
}

func ReadConfigFile() error {
	//for now statically assigning user id
	// ConfigObj.UserId = "6893814a3b3b86cffb0eaea1"
	// ConfigObj.NodeIPAddr = "100.81.30.122"
	configFile, err := os.OpenFile(APP_CONFIG_FILE_LOCATION, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("[ERROR] error opeing the config file ", APP_CONFIG_FILE_LOCATION)
		return err
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&ConfigObj)
	if err != nil {
		log.Println("[ERROR] error decoding the file to configObj")
		return err
	}
	configFile.Close()
	ValidateConfigFile()
	return nil
}

func ValidateConfigFile() error {
	//veryfyings public private keys
	if len(ConfigObj.PrivateKey) < 44 || len(ConfigObj.PublicKey) < 44 {
		privateKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return err
		}
		ConfigObj.PrivateKey, ConfigObj.PublicKey = privateKey.String(), privateKey.PublicKey().String()
	}
	// TODO: validate the public private key pair ( pretty simple )

	// verifying userID & nodeID & nodeIPAddrs
	// TODO: check authentication

	WriteConfigFile()
	//veryfing
	return nil
}
