package main

import (
	"flag"
	"net/http"
	"os"

	"bitbucket.org/cryptopatron/backend/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("changes test"))
}

func main() {
	servePath := flag.String("servePath", "../front-end/dist/dropcoin/", "Path to serve static files from")
	flag.Parse()

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Setup file serving from web app
	setupFileServer(router, *servePath)

	// Connect to DB


	// Setup REST API endpoints
	router.Get("/login", LoginHandler)
	router.Post("/auth/google/jwt", auth.GoogleAuthHandler)
	router.Post("/user/create", db.CreateUserHandler)

	// Our application will run on port 8080. Here we declare the port and pass in our router.
	http.ListenAndServe(":8008", router)

}

// Original source: https://github.com/go-chi/chi/issues/403
func setupFileServer(router *chi.Mux, root string) {
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Check if request 'r' is asking for a static file that does not exist,
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			// Serve the default index.html
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
			// w.Write([]byte("404 son"))
		} else {
			// Serve the file
			fs.ServeHTTP(w, r)
		}
	})
}
