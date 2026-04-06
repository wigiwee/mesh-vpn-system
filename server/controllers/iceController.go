package controllers

import (
	"encoding/json"
	"net/http"
	"server/config"
	"server/models"

	"github.com/gorilla/mux"
)

func RegisterIceCandidate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var candidateRegisterReq models.RegisterCandidateRequest
	json.NewDecoder(r.Body).Decode(&candidateRegisterReq)
	config.InMemoryCandidates[candidateRegisterReq.ConnectionIdentifier] =
		append(config.InMemoryCandidates[candidateRegisterReq.ConnectionIdentifier], candidateRegisterReq.Candidate)

	w.WriteHeader(http.StatusCreated)
}

func GetConnectionsCandidates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pathVars := mux.Vars(r)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		candidates []string
	}{
		candidates: config.InMemoryCandidates[models.ConnectionIdentifier{
			UserId:       pathVars["user_id"],
			LocalNodeId:  pathVars["local_node_id"],
			RemoteNodeId: pathVars["remote_node_id"],
		}],
	})
}

func RegisterCredentials(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerCredsReq models.RegisterCredentialsRequest
	json.NewDecoder(r.Body).Decode(&registerCredsReq)
	config.InMemoryCredentials[registerCredsReq.ConnectionIdentifier] = registerCredsReq.ICECreds
	w.WriteHeader(http.StatusCreated)
}

func GetCredentials(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pathVars := mux.Vars(r)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(config.InMemoryCredentials[models.ConnectionIdentifier{
		LocalNodeId:  pathVars["local_node_id"],
		RemoteNodeId: pathVars["remote_node_id"],
		UserId:       pathVars["user_id"],
	}])

}
