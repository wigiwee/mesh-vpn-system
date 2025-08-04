package main

const WG_CONFIG_FILE_LOCATION = "wg0.conf"

const SERVER_URL string = "http://localhost:4000"

const USER_ID string = "6888b831d2719d56bb8fdd7a"

var NODE_ID string

type AppState struct {
	IsLoggedIn bool
}
