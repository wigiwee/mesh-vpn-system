package controllers

import (
	"fmt"
	"net/http"
	"server/config"

	"github.com/gorilla/mux"
)

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	pathVars := mux.Vars(r)

	userId := pathVars["user_id"]
	// nodeId := pathVars["node_id"]

	ch := make(chan string)
	userChannels := config.UserNodesChannels[userId]
	userChannels = append(userChannels, ch)
	config.UserNodesChannels[userId] = userChannels

	defer func() {
		fmt.Println("ending connection")
		//remove the conn obj from the server
	}()
	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("Could not init http.Flusher")
	}

	for {
		select {
		case message := <-ch:
			fmt.Println("case message... sending message")
			fmt.Println(message)
			fmt.Fprintf(w, "data: %s\n\n", message)
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("Client closed connection")
			return
		}
	}
}
