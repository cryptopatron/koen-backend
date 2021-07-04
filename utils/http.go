package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

func DecodeJSON(w http.ResponseWriter, r *http.Request, schema interface{}, allowUnknownFields bool) error {

	if r.Body == nil {
		return errors.New("body cannot be empty")
	}

	decoder := json.NewDecoder(r.Body)
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
