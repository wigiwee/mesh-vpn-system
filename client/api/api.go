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
)

func RegisterNode(registerReq models.RegisterNodeRequest) (models.RegisterNodeResponse, error) {
	data, err := json.Marshal(registerReq)
	var registerNodeRes models.RegisterNodeResponse
	if err != nil {
		return registerNodeRes, err
	}
	resp, err := http.Post(config.SERVER_URL+"/api/registerNode", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return registerNodeRes, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return registerNodeRes, err
	}
	err = json.Unmarshal(body, &registerNodeRes)
	if err != nil {
		return registerNodeRes, err
	}
	log.Println("register node server response : " + registerNodeRes.NodeId + " " + registerNodeRes.IPAddress)
	// if resp.Status == strconv.Itoa(http.StatusOK) {
	// 	return models.RegisterNodeResponse{}, fmt.Errorf("%s", string(body))
	// }
	return registerNodeRes, nil
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
