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

	err := config.WriteWGConfig()
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

		//fetching node peers
		log.Println("[INFO] starting application loop ", config.PeerState)
		peers, err := api.GetPeers(config.ConfigObj.UserId, config.ConfigObj.NodeId)
		if err != nil {
			log.Println("[ERROR] error fetching peers ", err.Error())
		}
		log.Println("[DEBUG] received peers ", peers)

		//comparing newly fetched peers
		added, removed := helper.SyncPeers(peers)
		log.Println("[DEBUG] calcuated peer difference as added: ", added, " removed: ", removed)

		log.Println("[INFO] syncing the peers")
		// for _, peer := range added {
		// 	//this is where we generate the agent and establish the connection
		// 	agent := GetAgent()
		// 	id, pwd, _ := agent.GetLocalUserCredentials()
		// 	candidates := []string{}
		// 	agent.OnCandidate(func(c ice.Candidate) {
		// 		if c == nil {
		// 			return
		// 		}
		// 		log.Println("[INFO] found candidate " + c.String())
		// 		candidates = append(candidates, c.String())
		// 	})
		// 	agent.GatherCandidates()
		// 	time.Sleep(2 * time.Second)

		// 	api.RegisterIceCreds(models.ICECredsRegisterRequest{
		// 		LocalNodeId:  config.ConfigObj.NodeId,
		// 		RemoteNodeId: peer.NodeId,
		// 		UserId:       config.ConfigObj.UserId,
		// 		ICECreds: models.ICECreds{
		// 			ICEUfrag:   id,
		// 			ICEPwd:     pwd,
		// 			Candidates: candidates,
		// 		},
		// 	})

		// 	remoteCreds, err := api.FetchIceCreds(models.ICECredsFetchRequest{
		// 		LocalNodeId:  peer.NodeId,
		// 		RemoteNodeId: config.ConfigObj.NodeId,
		// 		UserId:       config.ConfigObj.UserId,
		// 	})
		// 	if err != nil {

		// 	}
		// 	agent.SetRemoteCredentials(id, pwd)
		// 	for _, candidate := range remoteCreds.Candidates {
		// 		candidate, err := ice.UnmarshalCandidate(candidate)
		// 		if err != nil {
		// 			fmt.Println("error decoding the candidate ", candidate, err)
		// 		}
		// 		agent.AddRemoteCandidate(candidate)
		// 	}
		// 	var conn net.Conn
		// 	if config.ConfigObj.NodeId > peer.NodeId {
		// 		conn, _ = agent.Dial(context.Background(), id, pwd)
		// 	} else {
		// 		conn, _ = agent.Accept(context.Background(), id, pwd)
		// 	}

		// 	config.AddPeer(peer)
		// 	config.PeerState[peer.PublicKey] = models.PeerState{
		// 		Peer:       peer,
		// 		Agent:      agent,
		// 		Conn:       &conn,
		// 		LocalUFrag: id,
		// 		LocalPwd:   pwd,
		// 		Connected:  true,
		// 	}
		// }

		for _, peer := range removed {
			config.RemovePeer(peer)
			delete(config.PeerState, peer.PublicKey)
		}
		time.Sleep(10 * time.Second)
	}
}
