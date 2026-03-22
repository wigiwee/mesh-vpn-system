package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
			Hostname:   node.Hostname,
			PublicKey:  node.PublicKey,
			IPAddress:  node.IPAddress,
			Endpoint:   node.Endpoint,
			NodeId:     node.Id.Hex(),
			ICEUfrag:   node.ICEUfrag,
			ICEPwd:     node.ICEPwd,
			Candidates: node.Candidates,
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(peers)
}

func UpdatePeer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

func UpdatePeersICECreds(w http.ResponseWriter, r *http.Request) {
	log.Println("got a update credentials request")
	w.Header().Set("Content-Type", "application/json")

	var iceCredsUpdateReq models.ICECredsUpdateRequest
	json.NewDecoder(r.Body).Decode(&iceCredsUpdateReq)
	log.Println("got the request object", iceCredsUpdateReq)
	err := services.UpdateICECreds(iceCredsUpdateReq)
	if err != nil {
		log.Print("error updating the nodes ice creds", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error updating the nodes ice creds" + err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
