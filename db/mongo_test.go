package db

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUserHandler(t *testing.T) {
	// Connect to DB
	var conn DBConn = &MongoInstance{Database: "koen_test", Collection: "users"}
	conn.Open()
	defer conn.Close()

	handler := HandleCreateUser(conn)

	t.Run("Bad request on empty request body", func(t *testing.T) {
		// Create a HTTP request with no JWT
		req, err := http.NewRequest("POST", "/user/create", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		handler.ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusBadRequest

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("HTTP 400 on random JSON request body", func(t *testing.T) {
		// Create a HTTP request with no JWT
		body := bytes.NewBuffer([]byte("bleh"))
		req, err := http.NewRequest("POST", "/user/create", body)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		handler.ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusBadRequest

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("HTTP 400 on incorrectly matching JSON fields", func(t *testing.T) {
		// Random fields which are dissimilar from the expected User creation fields
		var json string = `{"bleh":"fakeasstoken", "hello": 65}`
		body := strings.NewReader(json)
		req, err := http.NewRequest("POST", "/user/create", body)
		if err != nil {
			t.Fatal(err)
		}

		// TODO: Test to check if Content-Type is being checked for JSON
		// req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		handler.ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusBadRequest

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("HTTP 200 on correctly matching JSON request", func(t *testing.T) {
		// JSON which follow User key semantics in DB
		var json string = `{
			"pageName":"fakeasstoken",
			"name":"fakeasstoken",
			"email":"fakeasstoken",
			"metaMaskWalletPublicKey":"fakeasstoken",    
			"generatedMaticWalletPublicKey": "kuhgihjygyuh"
			}`
		body := strings.NewReader(json)
		req, err := http.NewRequest("POST", "/test", body)
		if err != nil {
			t.Fatal(err)
		}

		// TODO: Test to check if Content-Type is being checked for JSON
		// req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		handler.ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusOK

		if got != want {
			t.Log(rr)
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestHandleGetUser(t *testing.T) {
	// Connect to DB
	var conn DBConn = &MongoInstance{Database: "koen_test", Collection: "users"}
	conn.Open()
	defer conn.Close()

	handler := HandleGetUser(conn)

	t.Run("Find user even with just one field", func(t *testing.T) {
		// Random fields which are dissimilar from the expected User creation fields
		var json string = `{"email":"fakeasstoken"}`
		body := strings.NewReader(json)
		req, err := http.NewRequest("POST", "/test", body)
		if err != nil {
			t.Fatal(err)
		}

		// TODO: Test to check if Content-Type is being checked for JSON
		// req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		handler.ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("Find user on fully correct JSON request", func(t *testing.T) {
		// JSON which follow User key semantics in DB
		var json string = `{
			"profileName":"fakeasstoken",
			"userName":"fakeasstoken",
			"email":"fakeasstoken",
			"metaMaskWalletPublicKey":"fakeasstoken",    
			"autoWalletPublicKey": "kuhgihjygyuh"
			}`
		body := strings.NewReader(json)
		req, err := http.NewRequest("POST", "/test", body)
		if err != nil {
			t.Fatal(err)
		}

		// TODO: Test to check if Content-Type is being checked for JSON
		// req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		handler.ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusOK

		if got != want {
			t.Log(rr)
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
