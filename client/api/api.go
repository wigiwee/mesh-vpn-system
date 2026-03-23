package api

import (
	"bytes"
	"client/config"
	"client/models"
	"encoding/json"
	"errors"
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

func RegisterIceCreds(iceCredsRegisterRequest models.ICECredsRegisterRequest) error {
	data, err := json.Marshal(iceCredsRegisterRequest)
	if err != nil {
		log.Println("error encoding the iceCredsRegisterRequest")
		return err
	}
	resp, err := http.Post(config.SERVER_URL+"/api/registerCreds", "application/json", bytes.NewBuffer(data))
	log.Println("[INFO] hitting url " + config.SERVER_URL + "/api/registerCreds")
	if err != nil {
		log.Println("error register the creds")
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("got status code ", resp.StatusCode)
		return errors.New("got status code " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}

func FetchIceCreds(iceCredsFetchRequest models.ICECredsFetchRequest) (models.ICECreds, error) {
	var iceCreds models.ICECreds

	u, err := url.Parse(config.SERVER_URL + "/api/fetchCreds")
	if err != nil {
		return iceCreds, err
	}
	u.Scheme = "http"
	q := u.Query()
	q.Set("local_node_id", iceCredsFetchRequest.LocalNodeId)
	q.Set("remote_node_id", iceCredsFetchRequest.RemoteNodeId)
	q.Set("user_id", iceCredsFetchRequest.UserId)
	u.RawQuery = q.Encode()

	log.Println("hitting url " + u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		log.Println("error fetching the ice creds")
		return iceCreds, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&iceCreds); err != nil {
		log.Println("error decoding the response object into iceCreds obj")
		return iceCreds, err
	}
	if resp.StatusCode == 404 {
		return iceCreds, errors.New("ice agent at the peer not started yet")
	}

	return iceCreds, nil
}
