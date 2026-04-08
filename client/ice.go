package main

import (
	"client/api"
	"client/config"
	"client/models"
	"fmt"
	"log"

	"github.com/pion/ice/v2"
	"github.com/pion/stun"
)

func GetAgent(remoteNodeId string) *ice.Agent {

	agentConfig := &ice.AgentConfig{

		NetworkTypes: []ice.NetworkType{
			ice.NetworkTypeUDP4,
		},
		Urls: []*stun.URI{
			{
				Scheme: stun.SchemeTypeSTUN,
				Host:   "stun.l.google.com",
				Port:   19302,
			},
		},
	}
	agent, err := ice.NewAgent(agentConfig)
	if err != nil {
		log.Panic("[ERROR] error creating ice agent ", err)
	}

	agent.OnConnectionStateChange(func(state ice.ConnectionState) {
		log.Println("[INFO] ICE stagechange to " + state.String())
	})

	agent.OnCandidate(func(c ice.Candidate) {
		if c == nil {
			fmt.Println("candidate gathering finished")
			return
		}
		fmt.Println("booyah got a candidate ", c.String())
		err := api.AddCandidate(models.ConnectionIdentifier{
			UserId:       config.ConfigObj.UserId,
			LocalNodeId:  config.ConfigObj.NodeId,
			RemoteNodeId: remoteNodeId,
		}, c.Marshal())
		if err != nil {
			log.Println("[ERROR] adding the candiate", remoteNodeId)
		}
		fmt.Println("added candidate")
	})

	agent.GatherCandidates()

	return agent

}
