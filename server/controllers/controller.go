package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/models"
	"server/services"
)

// func getUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	json.NewEncoder(w).Encode("a;lksdjf")
// }

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

func GetPeers(w http.ResponseWriter, r *http.Request) {
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
