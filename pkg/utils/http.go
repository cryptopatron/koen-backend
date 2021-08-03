package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Respond(statusCode int, statusMsg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(statusMsg))
	}
}

// Helper function to return JSON responses
func RespondWithJSON(data interface{}, statusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		js, err := json.Marshal(data)
		// Make it easier to view JSON response in the Terminal
		js = append(js, '\n')
		if err != nil {
			http.Error(w, "The server encountered a problem,", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(js)
	}
}

func DecodeJSON(body io.Reader, schema interface{}, allowUnknownFields bool) error {

	if body == nil {
		return errors.New("body cannot be empty")
	}

	decoder := json.NewDecoder(body)
	if !allowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	err := decoder.Decode(&schema)

	if err != nil {
		log.Print("Decode error")
		log.Print(err)
		return errors.New("couldn't decode JSON")
	}

	return nil

}

type HttpTestCase struct {
	Handler http.HandlerFunc
	req     *http.Request
}

func (htc HttpTestCase) CheckReturnStatus(want int) (fn func(t *testing.T)) {
	return func(t *testing.T) {
		rr := httptest.NewRecorder()
		htc.Handler.ServeHTTP(rr, htc.req)

		got := rr.Code

		if got != want {
			t.Log(rr)
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func (htc *HttpTestCase) SetRequestBody(body io.Reader) {
	req, _ := http.NewRequest("POST", "/testCase", body)
	htc.req = req
}

func (htc *HttpTestCase) SetContext(key interface{}, val interface{}) {
	ctx := context.WithValue(htc.req.Context(), key, val)
	htc.req = htc.req.WithContext(ctx)
}
