package main

import (
	"client/api"
	"client/config"
	"client/helper"
	"client/models"
	"crypto/rand"
	"log"
	"os"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// TODO: change all the id fields type string->primitive.ObjectId
func main() {

	// TODO: check auth first
	config.ReadConfigFile()
	log.Println("[DEBUG] read the config : ", config.ConfigObj)
	if config.ConfigObj.UserId == "" {
		log.Println("[INFO] user logged out, logging in the user")
		//TODO: register/authenticate user
		config.ConfigObj.UserId = "6893814a3b3b86cffb0eaea1"
		log.Println("[INFO] user logged in")
	}
	if config.ConfigObj.NodeId == "" {
		log.Println("[INFO] node not registered, registering the node with the server")

		//fetching endpoint
		endpoint, err := GetPublicEndpoint()
		if err != nil {
			log.Println("[ERROR] STUN failed: " + err.Error())
		}
		config.ConfigObj.Endpoint = endpoint
		//fetching hostname
		hostname, err := os.Hostname()
		if err != nil {
			log.Println("[ERROR] couldn't fetch hostname" + err.Error())
			config.ConfigObj.Hostname = rand.Text()[:6]
		} else {
			config.ConfigObj.Hostname = hostname
		}
		//generating public private kyes
		privateKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			log.Panic("[ERROR] can't generate public private keys ", err.Error())
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
		if err != nil {
			log.Panic("[ERROR] error registering the node", err)
		}
		log.Println("[INFO] node registered with nodeID ", registerNodeRes.NodeId)
		config.ConfigObj.NodeId = registerNodeRes.NodeId
		config.ConfigObj.NodeIPAddr = registerNodeRes.IPAddress

		config.WriteConfigFile()
		log.Println("[INFO] config file written")
	}

	agent := StartICE()

	config.ICEUfrag, config.ICEPwd, _ = agent.GetLocalUserCredentials()
	log.Println("[DEBUG] got credentials as ", config.ICEUfrag, config.ICEPwd)

	id, err := primitive.ObjectIDFromHex(config.ConfigObj.NodeId)
	if err != nil {
		log.Panic("invalid node id ", config.ConfigObj.NodeId)
	}

	// TODO: implement channel here instead of just waiting for n seconds
	time.Sleep(3 * time.Second)
	log.Println("candidates ", config.Candidates)
	api.UpdateIceCreds(models.ICECredsUpdateRequest{
		Id:         id,
		ICEUfrag:   config.ICEUfrag,
		ICEPwd:     config.ICEPwd,
		Candidates: config.Candidates,
	})
	log.Println("[INFO] updated the ice creds to the server")

	err = config.WriteWGConfig()
	if err != nil {
		log.Printf("[ERROR] cannot write the %s.conf file %s", config.INTERFACE_NAME, err.Error())
	}

	cmd := exec.Command("sudo", "wg-quick", "up", config.WG_CONFIG_FILE_LOCATION)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("[ERROR]: %s\n", err)
	}
	log.Println("[INFO] brought up the wireguard interaface " + config.INTERFACE_NAME + "✅")

	//Daemon loop
	for {

		log.Println("[INFO] starting application loop ", config.Peers)
		peers, err := api.GetPeers(config.ConfigObj.UserId, config.ConfigObj.NodeId)
		if err != nil {
			log.Println("[ERROR] error fetching peers ", err.Error())
		}
		log.Println("[DEBUG] received peers ", peers)
		added, removed := helper.SyncPeers(peers)
		log.Println("[DEBUG] calcuated peer difference as added: ", added, " removed: ", removed)

		log.Println("[INFO] syncing the peers")
		for _, peer := range added {
			agent.SetRemoteCredentials(peer.ICEUfrag, peer.ICEPwd)
			config.AddPeer(peer)
			config.Peers[peer.PublicKey] = peer
		}
		for _, peer := range removed {
			config.RemovePeer(peer)
			delete(config.Peers, peer.PublicKey)
		}
		time.Sleep(10 * time.Second)
	}
}
