package config

import "server/models"

var (
	InMemoryCredentials map[models.ConnectionIdentifier]models.ICECreds = make(map[models.ConnectionIdentifier]models.ICECreds)
	InMemoryCandidates  map[models.ConnectionIdentifier][]string        = make(map[models.ConnectionIdentifier][]string)
)

const (
	Port int = 4000
)
