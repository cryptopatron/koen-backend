package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cryptopatron/koen-backend/pkg/utils"
)

func generateTestJWT() (JWT, error) {
	payload := generateTestPayload("Hello hash")
	// Payload to json string
	body, err := json.Marshal(payload)
	if err != nil {
		return JWT{}, err
	}

	req, err := http.NewRequest("POST", "/test", bytes.NewReader(body))
	if err != nil {
		return JWT{}, err
	}

	// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
	rr := httptest.NewRecorder()

	// handler satisfies the interface of http.Handler
	// So we can use its ServeHTTP to serve the rquest to it
	HandleWalletAuthentication().ServeHTTP(rr, req)

	// Buffer to JSON
	jwt := &JWT{}
	err = utils.DecodeJSON(rr.Body, jwt, false)
	if err != nil {
		return JWT{}, err
	}
	fmt.Print("JWT", jwt)
	return *jwt, nil
}

func TestHandleJWT(t *testing.T) {
	handler := HandleJWT(utils.Respond(http.StatusOK, ""))

	t.Run("Bad request on empty request body", func(t *testing.T) {
		// Create a HTTP request with no JWT
		req, err := http.NewRequest("POST", "/auth/google/jwt", nil)
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

	t.Run("Bad request with random request body", func(t *testing.T) {
		// Create a HTTP request with no JWT
		body := bytes.NewBuffer([]byte("bleh"))
		req, err := http.NewRequest("POST", "/auth/google/jwt", body)
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

	t.Run("Auth failed on random JWT", func(t *testing.T) {

		var jwt string = `{"idToken":"fakeasstoken"}`
		body := strings.NewReader(jwt)
		req, err := http.NewRequest("POST", "/auth/google/jwt", body)
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
		want := http.StatusUnauthorized

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("HTTP 200 on valid Wallet-based JWT", func(t *testing.T) {
		// Create a new JWT
		jwt, err := generateTestJWT()
		if err != nil {
			t.Error(err)
		}

		var body string = fmt.Sprintf(`{
				"idToken": "%s",
				"random": "random",
				"metaMaskWalletPublicKey": "",
				"generatedMaticWalletPublicKey": "kuhgihjygyuh"
			}`, jwt.IdToken)

		req, err := http.NewRequest("POST", "/test", strings.NewReader(body))
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

	// TODOL HTTP 200 on Google-based JWT
}
