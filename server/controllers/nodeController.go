package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/config"
	"server/models"
	"server/services"

	"github.com/gorilla/mux"
)

func RegisterNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerNodeReq models.RegisterNodeRequest
	json.NewDecoder(r.Body).Decode(&registerNodeReq)
	registerNodeRes, err := services.AddNode(registerNodeReq)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registerNodeRes)
}

func GetPeersOfUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	fmt.Println("received user_id: " + params["user_id"])
	nodes, err := services.FetchUserNodes(params["user_id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nodes)
}

func GetPeersOfNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(r.URL.String())
	query := r.URL.Query()
	log.Println("received node_id" + query.Get("node_id"))
	log.Println("received user_id" + query.Get("user_id"))
	nodeId := query.Get("node_id")
	userId := query.Get("user_id")
	nodes, err := services.FetchUserNodes(userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	//TODO: go through this one more time, nodeId & userId shouldn't be verified like this
	if len(nodeId) == 0 {
		log.Println("nodeId not found")
		w.Write([]byte("nodeId not found"))
		return
	}

	if len(userId) == 0 {
		log.Println("userId not found")
		w.Write([]byte("userId not found"))
		return
	}
	var peers []models.Peer = []models.Peer{}
	for _, node := range nodes {
		if node.Id.Hex() == nodeId {
			continue
		}
		peers = append(peers, models.Peer{
			Hostname:  node.Hostname,
			PublicKey: node.PublicKey,
			IPAddress: node.IPAddress,
			Endpoint:  node.Endpoint,
			NodeId:    node.Id.Hex(),
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(peers)
}

func UpdatePeer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

func RegisterICECreds(w http.ResponseWriter, r *http.Request) {
	log.Println("got a update credentials request")
	w.Header().Set("Content-Type", "application/json")

	var iceCredsRegisterReq models.ICECredsRegisterRequest
	json.NewDecoder(r.Body).Decode(&iceCredsRegisterReq)
	log.Println("got the request object", iceCredsRegisterReq)

	_, doesExist := config.InMemoryCredentials[iceCredsRegisterReq.UserId+iceCredsRegisterReq.LocalNodeId+iceCredsRegisterReq.RemoteNodeId]
	if doesExist {
		config.InMemoryCredentials[iceCredsRegisterReq.UserId+iceCredsRegisterReq.LocalNodeId+iceCredsRegisterReq.RemoteNodeId] = iceCredsRegisterReq.ICECreds
	}
	log.Println("successfully registered the creds in the memorydb")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func FetchICECreds(w http.ResponseWriter, r *http.Request) {
	log.Println("got a fetch icecreds req")
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	log.Println("received local_node_id", query.Get("local_node_id"))
	log.Println("received remote_node_id", query.Get("remote_node_id"))
	log.Println("received user_Id", query.Get("user_id"))
	localNodeId := query.Get("local_node_id")
	remoteNodeId := query.Get("remote_node_id")
	userId := query.Get("user_id")

	creds, doesExist := config.InMemoryCredentials[userId+localNodeId+remoteNodeId]
	if doesExist == false {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("no credentials found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(creds)
}
