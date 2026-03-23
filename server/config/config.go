package config

import "server/models"

var (
	InMemoryCredentials map[string]models.ICECreds = make(map[string]models.ICECreds)
)
