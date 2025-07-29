package routers

import (
	"server/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/api/registerUser", controllers.RegisterUser).Methods("POST")
	r.HandleFunc("/api/registerNode", controllers.RegisterNode).Methods("POST")
	r.HandleFunc("/api/getUserPeers", controllers.GetPeersOfUser).Methods("GET")
	r.HandleFunc("/api/getNodePeers", controllers.GetPeersOfNode).Methods("GET")
	// r.HandleFunc("/api/registerUser", RegisterUser())
	// r.HandleFunc("/api/registerUser", RegisterUser())

	return r
}
