package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"server/db"
	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(registerUserReq models.RegisterUserRequest) (string, error) {

	user := &models.User{
		Name:       registerUserReq.Name,
		Username:   registerUserReq.Username,
		Password:   registerUserReq.Password,
		NodesLimit: 100,
	}
	filter := bson.M{"username": registerUserReq.Username}
	cursor, err := db.UsersColl.Find(context.TODO(), filter)
	if err != nil {
		return "", err
	}
	var users []models.User
	cursor.All(context.TODO(), &users)
	if len(users) > 0 {
		return "", errors.New("user with provided userId exists")
	}
	ack, err := db.UsersColl.InsertOne(context.TODO(), user)
	if err != nil {
		return "nil", err
	}
	newUserId := ack.InsertedID.(primitive.ObjectID)

	log.Printf("added new user: %s -> %s \n", registerUserReq.Username, newUserId.Hex())
	return newUserId.Hex(), nil
}

func AddNode(newNodeReq models.RegisterNodeRequest) (string, error) {
	userPrimitiveId, err := primitive.ObjectIDFromHex(newNodeReq.UserId)
	if err != nil {
		log.Println(err)
		return "", err
	}
	newNode := &models.Node{
		AccessedBy: userPrimitiveId,
		Endpoint:   newNodeReq.Endpoint,
		IPAddress:  newNodeReq.IPAddress,
		Device:     newNodeReq.Device,
		PublicKey:  newNodeReq.PublicKey,
	}
	filter := bson.M{"endpoint": newNodeReq.Endpoint}
	cursor, err := db.NodesColl.Find(context.TODO(), filter)
	if err != nil {
		return "", err
	}
	var nodes []models.Node
	cursor.All(context.TODO(), &nodes)
	if len(nodes) > 0 {
		return "", errors.New("node with provided nodeId exists")
	}
	ack, err := db.NodesColl.InsertOne(context.TODO(), newNode)
	if err != nil {
		return "", err
	}
	newNodeId, _ := ack.InsertedID.(primitive.ObjectID)

	log.Printf("added new node: %s -> %s \n", newNode.Endpoint, newNodeId.Hex())
	return newNodeId.Hex(), nil
}

func FetchUserNodes(userId string) ([]models.Node, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"accessed_by": id}
	cursor, err := db.NodesColl.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var result []primitive.M
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		return nil, err
	}
	fmt.Println("result: ", result)
	var nodes []models.Node
	for _, m := range result {
		var node models.Node
		bsonBytes, err := bson.Marshal(m)
		if err != nil {
			return nil, err
		}
		err = bson.Unmarshal(bsonBytes, &node)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
