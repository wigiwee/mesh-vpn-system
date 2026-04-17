package config

import "server/models"

var (
	InMemoryCredentials map[models.ConnectionIdentifier]models.ICECreds = make(map[models.ConnectionIdentifier]models.ICECreds)
	InMemoryCandidates  map[models.ConnectionIdentifier][]string        = make(map[models.ConnectionIdentifier][]string)
	UserNodesChannels   map[string][]chan string                        = make(map[string][]chan string)
)

const (
	Port int = 4000
)
