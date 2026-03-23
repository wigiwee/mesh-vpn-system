package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"server/db"
	"server/helper"
	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddNode(newNodeReq models.RegisterNodeRequest) (models.RegisterNodeResponse, error) {
	log.Println("addnode method executing started")
	userPrimitiveId, err := primitive.ObjectIDFromHex(newNodeReq.UserId)
	if err != nil {
		log.Println(err)
		return models.RegisterNodeResponse{}, err
	}

	newNode := &models.Node{
		AccessedBy: userPrimitiveId,
		Endpoint:   newNodeReq.Endpoint,
		IPAddress:  helper.GenerateRandomIPAddr(),
		Device:     newNodeReq.Device,
		PublicKey:  newNodeReq.PublicKey,
		Hostname:   newNodeReq.Hostname,
	}

	filter := bson.M{"$or": []bson.M{
		{"endpoint": newNode.Endpoint},
		{"ip_address": newNode.IPAddress},
	}}
	cursor, err := db.NodesColl.Find(context.TODO(), filter)
	if err != nil {
		return models.RegisterNodeResponse{}, err
	}
	var nodes []models.Node
	cursor.All(context.TODO(), &nodes)
	if len(nodes) > 0 {
		return models.RegisterNodeResponse{}, errors.New("node with provided nodeId or IP address exists")
	}
	ack, err := db.NodesColl.InsertOne(context.TODO(), newNode)
	if err != nil {
		return models.RegisterNodeResponse{}, err
	}
	newNodeId, _ := ack.InsertedID.(primitive.ObjectID)

	log.Printf("added new node: %s -> %s \n", newNode.Endpoint, newNodeId.Hex())
	return models.RegisterNodeResponse{NodeId: newNodeId.Hex(), IPAddress: newNode.IPAddress}, nil
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

func UpdateNode(updateReq models.UpdateNodeRequest) error {
	filter := bson.M{"_id": updateReq.Id}
	update := bson.M{
		"$set": bson.M{
			"public_key": updateReq.PublicKey,
			"ip_address": updateReq.IPAddress,
			"endpoint":   updateReq.Endpoint,
			"hostname":   updateReq.Hostname,
		},
	}
	_, err := db.NodesColl.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
