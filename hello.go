package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Koen"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logincurl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.34.0/install.sh | bash"))
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", HomeHandler)
	router.Get("/login", LoginHandler)

	// Our application will run on port 8080. Here we declare the port and pass in our router.
	http.ListenAndServe(":8080", router)

}
