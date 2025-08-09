package models

type RegisterRequest struct {
	PublicKey string `json:"public_key" bson:"public_key"`
	IPAddress string `json:"ip_address" bson:"ip_address"`
	Endpoint  string `json:"endpoint"`
	Device    string `json:"device"`
	UserId    string `json:"user_id" bson:"user_id"`
}

type Node struct {
	Id         string `json:"id"`
	AccessedBy string `json:"accessed_by"`
	PublicKey  string `json:"public_key"`
	IPAddress  string `json:"ip_address"`
	Endpoint   string `json:"endpoint"`
	Device     string `json:"device"`
}
