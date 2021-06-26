package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/cryptopatron/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBConn interface {
	Open()
	Close()
	Create(entry interface{}) (interface{}, error)
}

type MongoInstance struct {
	client     *mongo.Client
	ctx        context.Context
	Database   string
	Collection string
}

func (m *MongoInstance) Open() {
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

func (m *MongoInstance) Close() {
	if err := m.client.Disconnect(m.ctx); err != nil {
		panic(err)
	}
}

func (m *MongoInstance) Create(entry interface{}) (interface{}, error) {
	doc, err := bson.Marshal(entry)
	if err != nil {
		return nil, err
	}
	collection := m.client.Database(m.Database).Collection(m.Collection)
	res, err := collection.InsertOne(m.ctx, doc)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	return id, err
}

func HandleCreateUser(db DBConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Body == nil {
			utils.Respond(http.StatusBadRequest, "Body is empty").ServeHTTP(w, r)
			return
		}

		defer r.Body.Close()

		type User struct {
			ProfileName    string `bson:"profileName"`
			AutoWalletAddr int    `bson:"autoWalletAddr"`
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		user := User{}
		err := decoder.Decode(&user)
		if err != nil {
			utils.Respond(http.StatusBadRequest, "Couldn't decode JSON").ServeHTTP(w, r)
			return
		}
		_, err = db.Create(user)
		if err != nil {
			utils.Respond(http.StatusInternalServerError, "Couldn't create new user!").ServeHTTP(w, r)
			return
		}

		utils.Respond(http.StatusOK, "")(w, r)
	}
}
