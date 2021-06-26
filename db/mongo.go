package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const DB_NAME = "koen"
const COLLECTION_NAME = "users"

type DBConn interface {
	open()
	close()
	create()
}

type MongoInstance struct {
	client *mongo.Client
	ctx context.Context
}

func (m *MongoInstance) open() {
	// Connect to MongoDB Atlas cloud DB
	uri := "mongodb+srv://koen:Rs19tpbHCpR7Lw7E@cluster0.znzdg.mongodb.net/koen?retryWrites=true&w=majority"
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	m.client = client
	m.ctx = ctx

}

func (m *MongoInstance) close() {
	if err := m.client.Disconnect(m.ctx); err != nil {
		panic(err)
	}
}

func (m *MongoInstance) create() {
	collection := m.client.Database(DB_NAME).Collection(COLLECTION_NAME)
	res, err := collection.InsertOne(m.ctx, bson.D{{"profileName", "test"}, {"value", 3.14159}})
	if err != nil {
		panic(err)
	}
	id := res.InsertedID
	fmt.Println(id)
}


type User struct {
	profileName string
	walletAddr int64
}

// func CreateUserHandler(w http.ResponseWriter, r *http.Request) {


// }

func main() {
	m := MongoInstance{}
	m.open()
	defer m.close()
	m.create()
	
}