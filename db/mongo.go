package db

import (
	"context"
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
	Read(query interface{}) (interface{}, error)
}

type MongoInstance struct {
	client     *mongo.Client
	ctx        context.Context
	Database   string
	Collection string
}
type User struct {
	Email                   string `bson:"email"` // Used for identifying Google users
	UserName                string `bson:"userName"`
	ProfileName             string `bson:"profileName"`
	AutoWalletPublicKey     string `bson:"autoWalletPublicKey"`
	MetaMaskWalletPublicKey string `bson:"metaMaskWalletPublicKey` // Used for identifying MetaMask users
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

func (m *MongoInstance) Read(query interface{}) (interface{}, error) {
	doc, err := bson.Marshal(query)
	if err != nil {
		return nil, err
	}
	var result bson.M
	collection := m.client.Database(m.Database).Collection(m.Collection)
	err = collection.FindOne(m.ctx, doc).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func HandleCreateUser(db DBConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &User{}
		err := utils.DecodeJSON(w, r, user)
		if err != nil {
			utils.Respond(http.StatusBadRequest, err.Error()).ServeHTTP(w, r)
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

func HandleGetUser(db DBConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &User{}
		err := utils.DecodeJSON(w, r, user)
		if err != nil {
			utils.Respond(http.StatusBadRequest, err.Error()).ServeHTTP(w, r)
			return
		}
		// Pass in an identifier struct
		result, err := db.Read(
			struct {
				Email string
			}{Email: user.Email})
		if err != nil {
			fmt.Print(err)
			http.NotFound(w, r)
			// utils.Respond(http.StatusOK, `{Couldn't find user!").ServeHTTP(w, r)
			return
		}
		fmt.Println(result)

		utils.Respond(http.StatusOK, "")(w, r)
	}
}
