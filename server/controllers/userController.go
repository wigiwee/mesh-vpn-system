package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/models"
	"server/services"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerUserReq models.RegisterUserRequest
	json.NewDecoder(r.Body).Decode(&registerUserReq)
	userId, err := services.AddUser(registerUserReq)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error registering node " + err.Error()))
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userId)
}

// TODO: write UpdateUser api
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updateReq models.UpdateNodeRequest
	json.NewDecoder(r.Body).Decode(&updateReq)

	err := services.UpdateNode(updateReq)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error updating node " + err.Error()))
	}
	w.WriteHeader(http.StatusOK)
}
