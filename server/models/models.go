package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Node struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AccessedBy primitive.ObjectID `json:"accessed_by" bson:"accessed_by"`
	PublicKey  string             `json:"public_key"`
	IPAddress  string             `json:"ip_address"`
	Endpoint   string             `json:"endpoint"`
	Device     string             `json:"device"`
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
	IPAddress string `json:"ip_address" bson:"ip_address"`
	Endpoint  string `json:"endpoint"`
	Device    string `json:"device"`
	UserId    string `json:"user_id" bson:"user_id"`
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
