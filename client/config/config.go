package config

const WG_CONFIG_FILE_LOCATION = "wg0.conf"

const SERVER_URL string = "http://localhost:4000"

const USER_ID string = "6893814a3b3b86cffb0eaea1"

const CONFIG_FILE_LOCATION = "local_files/"

var STUN_SERVERS []string = []string{
	"stun.l.google.com:19302",
	"stun1.l.google.com:19302",
	"stun2.l.google.com:19302",
	"stun3.l.google.com:19302",
	"stun4.l.google.com:19302"}

var NODE_ID string

type AppState struct {
	IsLoggedIn       bool
	IsNodeRegistered bool
	NodeIpAddr       string
}

func WriteConfig() {
}
