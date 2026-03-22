package main

import (
	"client/config"
	"log"

	"github.com/pion/ice/v2"
	"github.com/pion/stun"
)

func StartICE() *ice.Agent {

	agent := GetAgent()

	err := agent.GatherCandidates()
	if err != nil {
		log.Panic(err)
	}
	log.Println("ice gathering started")

	return agent
}

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
	agent.OnCandidate(func(c ice.Candidate) {
		if c == nil {
			return
		}
		log.Println("[INFO] found candidate " + c.String())
		config.Candidates = append(config.Candidates, c.String())
	})
	return agent

}
