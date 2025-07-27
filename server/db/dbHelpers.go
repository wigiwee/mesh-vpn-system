package db

import (
	"context"
	"errors"
	"log"
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
	cursor, err := usersColl.Find(context.TODO(), filter)
	if err != nil {
		return "", err
	}
	var users []models.User
	cursor.All(context.TODO(), &users)
	if len(users) > 0 {
		return "", errors.New("user with provided userId exists")
	}
	ack, err := usersColl.InsertOne(context.TODO(), user)
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
		IPAddress:  "100.100.100.100",
		Device:     newNodeReq.Device,
	}
	filter := bson.M{"endpoint": newNodeReq.Endpoint}
	cursor, err := nodesColl.Find(context.TODO(), filter)
	if err != nil {
		return "", err
	}
	var nodes []models.Node
	cursor.All(context.TODO(), &nodes)
	if len(nodes) > 0 {
		return "", errors.New("node with provided nodeId exists")
	}
	ack, err := nodesColl.InsertOne(context.TODO(), newNode)
	if err != nil {
		return "", err
	}
	newNodeId, _ := ack.InsertedID.(primitive.ObjectID)
	return newNodeId.Hex(), nil
}

func GetUsersNodes(userId string) ([]models.Node, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"accessed_by": id}
	cursor, err := nodesColl.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var result []primitive.M
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		return nil, err
	}
	marshalResult, err := bson.Marshal(result)
	if err != nil {
		return nil, err
	}
	var nodes []models.Node
	err = bson.Unmarshal(marshalResult, &nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
