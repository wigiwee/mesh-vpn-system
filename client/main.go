package main

import (
	"client/api"
	"client/config"
	"client/models"
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/pion/ice/v2"
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
		for _, peer := range peers {
			localConnIdentifier := models.ConnectionIdentifier{
				LocalNodeId:  config.ConfigObj.NodeId,
				RemoteNodeId: peer.NodeId,
				UserId:       config.ConfigObj.UserId}
			_, ifExist := config.PeerState[localConnIdentifier]
			if !ifExist {
				agent := GetAgent(peer.NodeId)
				id, pwd, _ := agent.GetLocalUserCredentials()
				creds := models.ICECreds{
					ICEUfrag: id,
					ICEPwd:   pwd,
				}
				err := api.RegisterIceCreds(localConnIdentifier, creds)
				if err != nil {
					log.Println("[ERROR] error registering the ice creds for remoteNode:", peer.NodeId)
				}
				config.PeerState[localConnIdentifier] = models.PeerState{
					Peer:       peer,
					Agent:      agent,
					Conn:       nil,
					LocalCreds: creds,
					// RemoteCreds:     models.ICECreds{ICEUfrag: "", ICEPwd: ""},
					ConnectedStatus:   false,
					IsRemoteConnected: false,
				}
			}
		}

		for localConnIdentifier, peerState := range config.PeerState {
			if peerState.RemoteCreds.ICEPwd == "" && peerState.RemoteCreds.ICEUfrag == "" {
				remoteConnIdentifier := models.ConnectionIdentifier{
					LocalNodeId:  localConnIdentifier.RemoteNodeId,
					RemoteNodeId: localConnIdentifier.LocalNodeId,
					UserId:       localConnIdentifier.UserId,
				}
				remoteCreds, err := api.GetIceCredentials(remoteConnIdentifier)
				if err != nil {
					log.Println("[ERROR] error fetching remote creds ", err)
					continue
				}
				candidates, err := api.GetCandidate(remoteConnIdentifier)
				if err != nil {
					log.Println("[ERROR] error fetching the candidates", err)
					continue
				}
				ps := config.PeerState[localConnIdentifier]
				ps.RemoteCreds = remoteCreds
				peerState.Agent.SetRemoteCredentials(remoteCreds.ICEUfrag, remoteCreds.ICEPwd)

				config.PeerState[localConnIdentifier] = ps
				//TODO: here candidates are added initially, it shouldn't be like this
				for _, candidateStr := range candidates {
					candidate, err := ice.UnmarshalCandidate(candidateStr)
					fmt.Println("candidate :: ", candidate.String())
					if err != nil {
						log.Println("[ERROR] error unmarshaling the candidates string ", candidateStr)
						continue
					}
					peerState.Agent.AddRemoteCandidate(candidate)
				}
				if localConnIdentifier.LocalNodeId > localConnIdentifier.RemoteNodeId {
					go dial(peerState.Agent, localConnIdentifier)
				} else {
					go accept(peerState.Agent, localConnIdentifier)
				}
			}
		}
		fmt.Println("i made it to here")

		for _, peerState := range config.PeerState {
			fmt.Println(peerState)
			fmt.Println(peerState.IsRemoteConnected)
			if peerState.IsRemoteConnected {
				fmt.Println("i shouldn't be here")
				config.AddPeer(peerState.Peer, peerState.Conn.RemoteAddr().String())
			}
		}
		fmt.Println("damn i made it to here")
		//TODO:ping the [ConnectedStatus: true] peers to check if they are still connected
		time.Sleep(10 * time.Second)
	}
}

func dial(agent *ice.Agent, connIdentifier models.ConnectionIdentifier) {
	log.Println("[INFO] dialing connection")
	ps := config.PeerState[connIdentifier]
	conn, err := agent.Dial(context.Background(), ps.RemoteCreds.ICEUfrag, ps.RemoteCreds.ICEPwd)
	if err != nil {
		log.Println(err)
	}

	ps.IsRemoteConnected = true
	ps.Conn = conn
	config.PeerState[connIdentifier] = ps
	log.Println("connection successful")
}

func accept(agent *ice.Agent, connIdentifier models.ConnectionIdentifier) {
	log.Println("[INFO] accepting connection")
	conn, _ := agent.Accept(context.Background(), "", "")
	peerState := config.PeerState[connIdentifier]
	peerState.IsRemoteConnected = true
	peerState.Conn = conn
	config.PeerState[connIdentifier] = peerState
	log.Println("connection successful")
}
