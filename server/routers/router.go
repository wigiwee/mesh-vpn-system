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
	r.HandleFunc("/api/registerNode", controllers.RegisterNode).Methods("POST")
	r.HandleFunc("/api/getUserPeers", controllers.GetPeersOfUser).Methods("GET")
	r.HandleFunc("/api/getNodePeers", controllers.GetPeersOfNode).Methods("GET")
	r.HandleFunc("/api/updateNode", controllers.UpdatePeer).Methods("PUT")
	r.HandleFunc("/api/updateNodeIceCreds", controllers.UpdatePeersICECreds).Methods("PUT")
	// r.HandleFunc("/api/registerUser", RegisterUser())
	// r.HandleFunc("/api/registerUser", RegisterUser())

	return r
}
