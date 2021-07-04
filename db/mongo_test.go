package db

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cryptopatron/koen-backend/auth"
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

	t.Run("HTTP 200 on partially correct JSON fields", func(t *testing.T) {
		// JSON which follow User key semantics but with some random fields
		var json string = `{
			"pageName":"fakeasstoken",
			"idToken":"bleh",
			"random": "random",
			"metaMaskWalletPublicKey":"fakeasstoken",    
			"generatedMaticWalletPublicKey": "kuhgihjygyuh"
			}`
		body := strings.NewReader(json)
		req, err := http.NewRequest("POST", "/test", body)
		ctx := context.WithValue(req.Context(), "userData",
			auth.GoogleClaims{
				Email: "test@koen.com", FirstName: "Koen", LastName: "San",
			})
		req = req.WithContext(ctx)
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
	t.Run("HTTP 200 on correctly matching JSON structure", func(t *testing.T) {
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
		ctx := context.WithValue(req.Context(), "userData",
			auth.GoogleClaims{
				Email: "test@koen.com", FirstName: "Koen", LastName: "San",
			})
		req = req.WithContext(ctx)
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

	t.Run("HTTP 200 on finding user with known email", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/test", nil)
		ctx := context.WithValue(req.Context(), "userData", auth.GoogleClaims{Email: "fakeasstoken"})
		req = req.WithContext(ctx)
		if err != nil {
			t.Fatal(err)
		}

		// TODO: test for HEader status codes
		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusOK

		if got != want {
			t.Log(rr)
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("HTTP 404 on sending empty google claim", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/test", nil)
		ctx := context.WithValue(req.Context(), "userData", auth.GoogleClaims{})
		req = req.WithContext(ctx)
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
		want := http.StatusNotFound

		if got != want {
			t.Log(rr)
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
