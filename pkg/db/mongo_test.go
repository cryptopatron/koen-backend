package db

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/cryptopatron/koen-backend/pkg/auth"
	"github.com/cryptopatron/koen-backend/pkg/utils"
)

func TestCreateUserHandler(t *testing.T) {
	// Connect to DB
	var conn DBConn = &MongoInstance{Database: "koen_test", Collection: "users"}
	conn.Open()
	defer conn.Close()

	htc := &utils.HttpTestCase{
		Handler: HandleCreateUser(conn),
	}

	htc.SetRequestBody(nil)
	t.Run("HTTP 400 on empty request body", htc.CheckReturnStatus(http.StatusBadRequest))

	htc.SetRequestBody(bytes.NewBuffer([]byte("bleh")))
	t.Run("HTTP 400 on random request body", htc.CheckReturnStatus(http.StatusBadRequest))

	// JSON which follow User key semantics but with some random fields
	var extraJson string = `{
		"pageName": "fakeasstoken",
		"idToken": "bleh",
		"random": "random",
		"metaMaskWalletPublicKey":"",    
		"generatedMaticWalletPublicKey": "kuhgihjygyuh",
		"random4": "random",
		"story": "of my life"
		}`
	htc.SetRequestBody(strings.NewReader(extraJson))
	claim := auth.Claims{
		GoogleClaims: auth.GoogleClaims{
			Email: "test@koen.com", FirstName: "Koen", LastName: "San",
		},
	}
	htc.SetContext("userData", claim)
	t.Run("HTTP 200 on correct JSON with extra fields", htc.CheckReturnStatus(http.StatusOK))

	var correctJson string = `{
		"pageName":"fakeasstoken",
		"name":"fakeasstoken",
		"email":"fakeasstoken",
		"metaMaskWalletPublicKey":"fakeasstoken",    
		"generatedMaticWalletPublicKey": "kuhgihjygyuh"
		}`
	htc.SetRequestBody(strings.NewReader(correctJson))
	htc.SetContext("userData", claim)
	t.Run("HTTP 200 on correctly matching JSON", htc.CheckReturnStatus(http.StatusOK))
}

func TestHandleGetUser(t *testing.T) {
	// Connect to DB
	var conn DBConn = &MongoInstance{Database: "koen_test", Collection: "users"}
	conn.Open()
	defer conn.Close()

	htc := &utils.HttpTestCase{
		Handler: HandleGetUser(conn),
	}

	htc.SetRequestBody(nil)
	htc.SetContext("userData", auth.Claims{
		GoogleClaims: auth.GoogleClaims{Email: "fakeasstoken"},
	})
	t.Run("HTTP 200 on getting user with GoogleClaims email", htc.CheckReturnStatus(http.StatusOK))
	// TODO: Also check body for user details

	// Need to test better
	htc.SetRequestBody(nil)
	htc.SetContext("userData", auth.Claims{})
	t.Run("HTTP 200 / Empty response on sending empty google claim", htc.CheckReturnStatus(http.StatusOK))

	htc.SetRequestBody(nil)
	htc.SetContext("userData", auth.Claims{
		WalletClaims: auth.WalletClaims{WalletPublicAddress: "fakeasstoken"},
	})
	t.Run("HTTP 200 on getting user with WalletClaims key", htc.CheckReturnStatus(http.StatusOK))
}
