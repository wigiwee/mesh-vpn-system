package api

import (
	"bytes"
	"client/config"
	"client/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func RegisterNode(registerReq models.RegisterRequest) (string, error) {
	data, err := json.Marshal(registerReq)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(config.SERVER_URL+"/api/registerNode", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Println("register node server response : " + string(body))
	if resp.Status == strconv.Itoa(http.StatusOK) {
		return "", fmt.Errorf(string(body))
	}
	return string(body), nil
}

func GetPeers(userId, nodeId string) ([]models.Node, error) {

	u, err := url.Parse(config.SERVER_URL + "/api/getNodePeers")
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
	var peers []models.Node
	if err := json.NewDecoder(resp.Body).Decode(&peers); err != nil {
		return nil, err
	}

	return peers, nil
}
