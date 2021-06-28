package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func Respond(statusCode int, statusMsg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(statusMsg))
	}
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, schema interface{}) error {

	if r.Body == nil {
		return errors.New("body cannot be empty")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&schema)
	if err != nil {
		return errors.New("couldn't decode JSON")
	}

	return nil

}
