package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.3.2"
const dbName = "CloadRoute"
const users = "users"
const nodes = "nodes"

var UsersColl *mongo.Collection
var NodesColl *mongo.Collection

func init() {
	clientOption := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal()
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Ping failed:", err)
	}
	log.Println("Mongodb connected")

	NodesColl = client.Database(dbName).Collection(nodes)
	UsersColl = client.Database(dbName).Collection(users)
}
