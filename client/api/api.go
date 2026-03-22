package api

import (
	"bytes"
	"client/config"
	"client/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func RegisterNode(registerReq models.RegisterNodeRequest) (models.RegisterNodeResponse, error) {
	data, err := json.Marshal(registerReq)
	var registerNodeRes models.RegisterNodeResponse
	if err != nil {
		return registerNodeRes, err
	}
	resp, err := http.Post(config.SERVER_URL+"/api/registerNode", "application/json", bytes.NewBuffer(data))
	log.Println("[INFO] hitting url : " + config.SERVER_URL + "/api/registerNode")
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
	log.Println("[DEBUG] register node server response : " + registerNodeRes.NodeId + " " + registerNodeRes.IPAddress)

	return registerNodeRes, nil
}

func GetPeers(userId, nodeId string) ([]models.Peer, error) {

	u, err := url.Parse(config.SERVER_URL + "/api/getNodePeers")
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "http"
	q := u.Query()
	q.Set("user_id", userId)
	q.Set("node_id", nodeId)
	u.RawQuery = q.Encode()

	log.Println("[INFO] hitting url : " + u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var peers []models.Peer
	if err := json.NewDecoder(resp.Body).Decode(&peers); err != nil {
		log.Println("[ERROR] error decoding response body ", err.Error())
		return nil, err
	}
	return peers, nil
}

func UpdateIceCreds(iceUpdateReq models.ICECredsUpdateRequest) error {
	data, err := json.Marshal(iceUpdateReq)
	if err != nil {
		log.Println("[ERROR] error encoding the request object ", err.Error())
		return err
	}
	req, err := http.NewRequest(http.MethodPut, config.SERVER_URL+"/api/updateNodeIceCreds", bytes.NewBuffer(data))
	log.Println("hitting req ", req.URL.String())

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] error updating the ice creds ", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		log.Println("[ERROR] something went wrong received status code " + strconv.Itoa(resp.StatusCode) + "body " + string(body))
		return err
	}

	return nil
}
