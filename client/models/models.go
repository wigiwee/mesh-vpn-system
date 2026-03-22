package models

import "go.mongodb.org/mongo-driver/bson/primitive"

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

	ICEUfrag   string   `json:"ice_ufrag" bson:"ice_ufrag"`
	ICEPwd     string   `json:"ice_pwd" bson:"ice_pwd"`
	Candidates []string `json:"candidates" bson:"candidates"`
}

type Peer struct {
	Hostname  string `json:"hostname"`
	PublicKey string `json:"public_key" bson:"public_key"`
	IPAddress string `json:"ip_address" bson:"ip_address"`
	Endpoint  string `json:"endpoint"`
	NodeId    string `json:"node_id" bson:"node_id"`

	ICEUfrag   string   `json:"ice_ufrag" bson:"ice_ufrag"`
	ICEPwd     string   `json:"ice_pwd" bson:"ice_pwd"`
	Candidates []string `json:"candidates" bson:"candidates"`
}

type ICECredsUpdateRequest struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ICEUfrag   string             `json:"ice_ufrag" bson:"ice_ufrag"`
	ICEPwd     string             `json:"ice_pwd" bson:"ice_pwd"`
	Candidates []string           `json:"candidates" bson:"candidates"`
}
