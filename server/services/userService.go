package services

import (
	"context"
	"errors"
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
