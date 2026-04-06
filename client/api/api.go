package api

import (
	"bytes"
	"client/config"
	"client/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func RegisterNode(registerReq models.RegisterNodeRequest) (models.RegisterNodeResponse, error) {
	data, err := json.Marshal(registerReq)
	var registerNodeRes models.RegisterNodeResponse
	if err != nil {
		return registerNodeRes, err
	}
	resp, err := http.Post(config.SERVER_URL+"/api/node", "application/json", bytes.NewBuffer(data))
	log.Println("[INFO] hitting url : " + config.SERVER_URL + "/api/node")
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

	url := config.SERVER_URL + fmt.Sprintf("/api/peer/%s/%s", userId, nodeId)
	log.Println("[INFO] hitting url : " + url)
	resp, err := http.Get(url)
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

func RegisterIceCreds(ConnectionIdentifier models.ConnectionIdentifier, iceCreds models.ICECreds) error {
	data, err := json.Marshal(models.RegisterCredentialsRequest{
		ConnectionIdentifier: ConnectionIdentifier,
		ICECreds:             iceCreds,
	})
	if err != nil {
		log.Println("[ERROR] error encoding the iceCredsRegisterRequest")
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/ice/candidate", config.SERVER_URL), "application/json", bytes.NewBuffer(data))
	log.Println("[INFO] hitting url " + fmt.Sprintf("%s/api/ice/candidate", config.SERVER_URL))
	if err != nil {
		log.Println("[ERROR] error register the creds")
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Println("got status code ", resp.StatusCode)
		return errors.New("got status code " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}

func FetchIceCreds(connectionIdentifier models.ConnectionIdentifier) (models.ICECreds, error) {
	var iceCreds models.ICECreds

	resp, err := http.Get(fmt.Sprintf("%s/api/ice/candidate/%s/%s/%s",
		config.SERVER_URL,
		connectionIdentifier.UserId,
		connectionIdentifier.LocalNodeId,
		connectionIdentifier.RemoteNodeId))
	log.Println("[INFO] hitting url to fetch ice creds")
	if err != nil {
		log.Println("[ERROR] error fetching the ice creds")
		return iceCreds, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&iceCreds); err != nil {
		log.Println("[ERROR] error decoding the response object into iceCreds obj")
		return iceCreds, err
	}
	if resp.StatusCode != http.StatusOK {
		return iceCreds, errors.New("[ERROR] ice agent at the peer not started yet")
	}
	return iceCreds, nil
}
