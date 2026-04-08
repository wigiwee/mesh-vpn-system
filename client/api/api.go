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
	url := fmt.Sprintf("%s/api/ice/credentials", config.SERVER_URL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	log.Println("[INFO] hitting url " + url)
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

func GetIceCredentials(connectionIdentifier models.ConnectionIdentifier) (models.ICECreds, error) {
	var iceCreds models.ICECreds

	resp, err := http.Get(fmt.Sprintf("%s/api/ice/credentials/%s/%s/%s",
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

	err = json.NewDecoder(resp.Body).Decode(&iceCreds)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		log.Println("[ERROR] error decoding the response object into iceCreds obj [OR] remote creds not available yet")
		return models.ICECreds{ICEUfrag: "", ICEPwd: ""}, err
	}
	//this if never gets executed because 404 throws err
	return iceCreds, nil
}

func AddCandidate(connectionIdentifier models.ConnectionIdentifier, candiate string) error {
	reqBody, err := json.Marshal(models.RegisterCandidateRequest{
		ConnectionIdentifier: connectionIdentifier,
		Candidate:            candiate,
	})
	resp, err := http.Post(fmt.Sprintf("%s/api/ice/candidate", config.SERVER_URL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("[ERROR] error registering the candidate", err)
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		log.Println("[ERROR] received status code " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}

func GetCandidate(connectionIdentifier models.ConnectionIdentifier) ([]string, error) {
	url := fmt.Sprintf("%s/api/ice/candidate/%s/%s/%s",
		config.SERVER_URL,
		connectionIdentifier.UserId,
		connectionIdentifier.LocalNodeId,
		connectionIdentifier.RemoteNodeId)
	resp, err := http.Get(url)

	log.Println("[INFO] hitting url ", url)
	if err != nil {
		log.Println("[ERROR] error hitting the url")
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("candidates not found")
	}
	var respObj struct {
		candidates []string
	}
	json.NewDecoder(resp.Body).Decode(&respObj)
	return respObj.candidates, nil
}
