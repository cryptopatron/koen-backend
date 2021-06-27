package utils

import "net/http"

func Respond(statusCode int, statusMsg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(statusMsg))
	}
}
