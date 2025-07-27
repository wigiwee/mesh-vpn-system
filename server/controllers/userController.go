package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"server/models"
)

// func getUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	json.NewEncoder(w).Encode("a;lksdjf")
// }

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerUserReq models.RegisterUserRequest
	json.NewDecoder(r.Body).Decode(&registerUserReq)
	userId, err := db.AddUser(registerUserReq)
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
	newNodeId, err := db.AddNode(registerNodeReq)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
	json.NewEncoder(w).Encode(newNodeId)
}
