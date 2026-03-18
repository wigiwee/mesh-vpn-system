package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Node struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	//TODO: rename the AccessedBy variable to UserId
	AccessedBy primitive.ObjectID `json:"accessed_by" bson:"accessed_by,omitempty"`
	PublicKey  string             `json:"public_key"`
	IPAddress  string             `json:"ip_address"`
	Endpoint   string             `json:"endpoint"`
	Device     string             `json:"device"`
	Hostname   string             `json:"hostname"`
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
