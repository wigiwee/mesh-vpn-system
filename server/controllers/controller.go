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

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerUserReq models.RegisterUserRequest
	json.NewDecoder(r.Body).Decode(&registerUserReq)
	userId, err := services.AddUser(registerUserReq)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
	json.NewEncoder(w).Encode(userId)
}

func RegisterNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerNodeReq models.RegisterNodeRequest
	json.NewDecoder(r.Body).Decode(&registerNodeReq)
	newNodeId, err := services.AddNode(registerNodeReq)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
	json.NewEncoder(w).Encode(newNodeId)
}

func GetPeersOfUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	fmt.Println("received user_id: " + params["user_id"])
	nodes, err := services.FetchUserNodes(params["user_id"])
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
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

	for idx, node := range nodes {
		if node.Id.Hex() == nodeId {
			nodes = append(nodes[:idx], nodes[idx+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(nodes)
}
