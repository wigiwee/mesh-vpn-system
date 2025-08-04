package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func registerNode(registerReq RegisterRequest) (string, error) {
	data, err := json.Marshal(registerReq)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(SERVER_URL+"/api/registerNode", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Println("register node server response : " + string(body))
	//TODO: implement status codes from server side and send appropriate responce accordingly
	return string(body), nil
}

func getPeers(userId, nodeId string) ([]Node, error) {

	u, err := url.Parse(SERVER_URL + "/api/getNodePeers")
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "http"
	q := u.Query()
	q.Set("user_id", userId)
	q.Set("node_id", nodeId)
	u.RawQuery = q.Encode()
	fmt.Println(u)

	log.Println("hitting url : " + u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var peers []Node
	if err := json.NewDecoder(resp.Body).Decode(&peers); err != nil {
		return nil, err
	}

	return peers, nil
}
