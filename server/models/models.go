package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Node struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AccessedBy primitive.ObjectID `json:"accessed_by" bson:"accessed_by,omitempty"` //TODO: rename the AccessedBy variable to UserId
	PublicKey  string             `json:"public_key" bson:"public_key"`
	IPAddress  string             `json:"ip_address" bson:"ip_address"`
	Endpoint   string             `json:"endpoint" bson:"endpoint"`
	Device     string             `json:"device" bson:"device"`
	Hostname   string             `json:"hostname" bson:"hostname"`
}

type UpdateNodeRequest struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PublicKey string             `json:"public_key"`
	IPAddress string             `json:"ip_address"`
	Endpoint  string             `json:"endpoint"`
	Device    string             `json:"device"`
	Hostname  string             `json:"hostname"`
}

type User struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username   string             `json:"username"`
	Name       string             `json:"name"`
	NodesLimit int                `json:"nodes_limit"`
	Password   string             `json:"password"`
}

type RegisterNodeRequest struct {
	PublicKey string `json:"public_key" bson:"public_key"`
	Endpoint  string `json:"endpoint"`
	Device    string `json:"device"`
	UserId    string `json:"user_id" bson:"user_id"`
	Hostname  string `json:"hostname" bson:"hostname"`
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type RegisterNodeResponse struct {
	IPAddress string `json:"ip_address" bson:"ip_address"`
	NodeId    string `json:"node_id" bson:"node_id"`
}

type Peer struct {
	Hostname  string `json:"hostname"`
	PublicKey string `json:"public_key" bson:"public_key"`
	IPAddress string `json:"ip_address" bson:"ip_address"`
	NodeId    string `json:"node_id" bson:"node_id"`
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

type ICECreds struct {
	ICEUfrag string `json:"ice_ufrag" bson:"ice_ufrag"`
	ICEPwd   string `json:"ice_pwd" bson:"ice_pwd"`
}

type RegisterCredentialsRequest struct {
	ConnectionIdentifier ConnectionIdentifier `json:"connection_identifier" bson:"connection_identifier"`
	ICECreds             ICECreds             `json:"ice_creds" bson:"ice_creds"`
}
