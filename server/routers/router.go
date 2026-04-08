package routers

import (
	"server/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/api/registerUser", controllers.RegisterUser).Methods("POST")

	// Nodes routes
	r.HandleFunc("/api/node", controllers.RegisterNode).Methods("POST")

	//peer routes
	r.HandleFunc("/api/peer/{user_id}", controllers.GetPeersOfUser).Methods("GET")
	r.HandleFunc("/api/peer/{user_id}/{node_id}", controllers.GetPeersOfNode).Methods("GET")

	//ice routes
	r.HandleFunc("/api/ice/credentials", controllers.RegisterCredentials).Methods("POST")
	r.HandleFunc("/api/ice/credentials/{user_id}/{local_node_id}/{remote_node_id}", controllers.GetCredentials).Methods("GET")
	r.HandleFunc("/api/ice/candidate", controllers.RegisterIceCandidate).Methods("POST")
	r.HandleFunc("/api/ice/candidate/{user_id}/{local_node_id}/{remote_node_id}", controllers.GetConnectionsCandidates).Methods("GET")

	// auth routes
	// r.HandleFunc("/api/registerUser", RegisterUser())
	// r.HandleFunc("/api/registerUser", RegisterUser())

	return r
}
