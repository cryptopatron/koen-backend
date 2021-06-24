package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGoogleAuthHandler(t *testing.T) {
	handler := http.HandlerFunc(GoogleAuthHandler)

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
}
