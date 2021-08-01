package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/cryptopatron/koen-backend/pkg/auth"
	"github.com/cryptopatron/koen-backend/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const API_PREFIX = "/api/v1"

func main() {
	servePath := flag.String("servePath", "../front-end/dist/dropcoin/", "Path to serve static files from")
	flag.Parse()

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Setup file serving from web app
	setupFileServer(router, *servePath)

	// Connect to DB
	var conn db.DBConn = &db.MongoInstance{Database: "koen_test", Collection: "users"}
	conn.Open()
	defer conn.Close()

	// TODO: Write tests for server endpoints
	// Setup REST API endpoints
	router.Route(API_PREFIX, func(r chi.Router) {
		r.Post("/auth/wallet", auth.HandleWalletAuthentication())

		// Protected routes
		r.Group(func(r chi.Router) {
			// Setup auth middleware
			r.Use(auth.HandleJWT)
			// Does not depend on auth mode
			r.Post("/users/create", db.HandleCreateUser(conn))

			// This depends on JWT context field: Email
			r.Post("/google/users/get", db.HandleGetUser(conn))

		})

		// Public routes
		r.Get("/users/pageName/{pageName}", func(w http.ResponseWriter, r *http.Request) {
			if pageName := chi.URLParam(r, "pageName"); pageName != "" {
				db.GetUser(conn, map[string]string{"pageName": pageName}).ServeHTTP(w, r)
			}
		})
	})

	// Our application will run on port 8080. Here we declare the port and pass in our router.
	fmt.Println("Running GO backend on port 8008")
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
