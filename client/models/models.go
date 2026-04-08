package models

import (
	"github.com/pion/ice/v2"
)

type RegisterNodeRequest struct {
	PublicKey string `json:"public_key" bson:"public_key"`
	IPAddress string `json:"ip_address" bson:"ip_address"`
	Endpoint  string `json:"endpoint"`
	Device    string `json:"device"`
	UserId    string `json:"user_id" bson:"user_id"`
	Hostname  string `json:"hostname" bson:"hostname"`
}

type RegisterNodeResponse struct {
	IPAddress string `json:"ip_address" bson:"ip_address"`
	NodeId    string `json:"node_id" bson:"node_id"`
}

type Node struct {
	Id         string `json:"id"`
	AccessedBy string `json:"accessed_by"`
	PublicKey  string `json:"public_key"`
	IPAddress  string `json:"ip_address"`
	Endpoint   string `json:"endpoint"`
	Device     string `json:"device"`
	Hostname   string `json:"hostname"`
}

type Peer struct {
	Hostname  string `json:"hostname"`
	PublicKey string `json:"public_key" bson:"public_key"`
	IPAddress string `json:"ip_address" bson:"ip_address"`
	NodeId    string `json:"node_id" bson:"node_id"`
}

type ICECreds struct {
	ICEUfrag string `json:"ice_ufrag" bson:"ice_ufrag"`
	ICEPwd   string `json:"ice_pwd" bson:"ice_pwd"`
}
type RegisterCredentialsRequest struct {
	ConnectionIdentifier ConnectionIdentifier `json:"connection_identifier" bson:"connection_identifier"`
	ICECreds             ICECreds             `json:"ice_creds" bson:"ice_creds"`
}

type RegisterCandidateRequest struct {
	ConnectionIdentifier ConnectionIdentifier `json:"connection_identifier" bson:"connection_identifier"`
	Candidate            string               `json:"candidate" bson:"candidate"`
}

type ConnectionIdentifier struct {
	LocalNodeId  string `json:"local_node_id" bson:"local_node_id"`
	RemoteNodeId string `json:"remote_node_id" bson:"remote_node_id"`
	UserId       string `json:"user_id" bson:"user_id"`
}

type PeerState struct {
	Peer Peer

	Agent *ice.Agent
	Conn  *ice.Conn

	LocalCreds  ICECreds
	RemoteCreds ICECreds

	//TODO: make the connected variable a enum, the states will be  connected/relayed/disconnected
	ConnectedStatus bool
	// PeerStatus        bool
	IsRemoteConnected bool
}
