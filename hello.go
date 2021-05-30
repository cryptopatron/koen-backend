package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Koen"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	http.Handle("/", router)

	// Our application will run on port 8080. Here we declare the port and pass in our router.
	http.ListenAndServe(":8080", router)

}
