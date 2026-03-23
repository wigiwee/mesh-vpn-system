package main

import (
	"log"

	"github.com/pion/ice/v2"
	"github.com/pion/stun"
)

func GetAgent() *ice.Agent {

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

	return agent

}
