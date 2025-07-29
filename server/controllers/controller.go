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
	fmt.Println(params["user_id"])
	nodes, err := services.FetchUserNodes(params["user_id"])
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
	json.NewEncoder(w).Encode(nodes)
}

func GetPeersOfNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	nodes, err := services.FetchUserNodes(params["user_id"])
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
	if len(params["node_id"]) == 0 {
		log.Println("nodeId not found")
		w.Write([]byte("nodeId not found"))
		return
	}
	for idx, node := range nodes {
		if node.Id.Hex() == params["node_id"] {
			nodes = append(nodes[:idx], nodes[idx+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(nodes)
}
