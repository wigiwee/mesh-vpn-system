package models

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
	Endpoint  string `json:"endpoint"`
	NodeId    string `json:"node_id" bson:"node_id"`
}
