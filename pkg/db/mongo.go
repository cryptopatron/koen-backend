package db

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cryptopatron/koen-backend/pkg/auth"
	"github.com/cryptopatron/koen-backend/pkg/utils"
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
	Email                         string `bson:"email" json:"email"` // Used for identifying Google users
	Name                          string `bson:"name" json:"name"`
	PageName                      string `bson:"pageName" json:"pageName"`
	GeneratedMaticWalletPublicKey string `bson:"generatedMaticWalletPublicKey" json:"generatedMaticWalletPublicKey"`
	MetaMaskWalletPublicKey       string `bson:"metaMaskWalletPublicKey json:"metaMaskWalletPublicKey"` // Used for identifying MetaMask users
}

func (m *MongoInstance) Open() {
	// Connect to MongoDB Atlas cloud DB
	uri := "mongodb+srv://koen:Rs19tpbHCpR7Lw7E@cluster0.znzdg.mongodb.net/koen?retryWrites=true&w=majority"
	ctx := context.Background()
	// defer cancel()
	client, err := mongo.Connect(m.ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	m.ctx = ctx
	m.client = client

	// Ping the primary
	if err := m.client.Ping(m.ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged MongoDB")

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
		err := utils.DecodeJSON(r.Body, user, true)
		if err != nil {
			utils.Respond(http.StatusBadRequest, err.Error()).ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		userData, ok := ctx.Value("userData").(auth.GoogleClaims)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		user.Name = userData.FirstName + " " + userData.LastName
		user.Email = userData.Email

		_, err = db.Create(user)
		if err != nil {
			utils.Respond(http.StatusInternalServerError, "Couldn't create new user!").ServeHTTP(w, r)
			return
		}
		utils.RespondWithJSON(user, http.StatusOK)(w, r)
	}
}

func HandleGetUser(db DBConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userData, ok := ctx.Value("userData").(auth.GoogleClaims)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
		// Pass in an identifier struct
		result, err := db.Read(
			struct {
				Email string
			}{Email: userData.Email})
		if err != nil {
			fmt.Print(err)
			utils.RespondWithJSON(struct{}{}, http.StatusOK).ServeHTTP(w, r)
			return
		}
		utils.RespondWithJSON(result, http.StatusOK)(w, r)
	}
}

func GetUser(db DBConn, searchQuery interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := db.Read(searchQuery)
		if err != nil {
			fmt.Print(err)
			utils.RespondWithJSON(struct{}{}, http.StatusOK).ServeHTTP(w, r)
			return
		}
		utils.RespondWithJSON(result, http.StatusOK)(w, r)
	}
}
