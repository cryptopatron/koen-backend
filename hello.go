package main

import (
	"net/http"
	"os"

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

	webAppPath := "../front-end/dist/dropcoin/"
	setupFileServer(router, webAppPath)

	// Our application will run on port 8080. Here we declare the port and pass in our router.
	http.ListenAndServe(":8080", router)

}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// func FileServer(r chi.Router, path string, root http.FileSystem) {
// 	if strings.ContainsAny(path, "{}*") {
// 		panic("FileServer does not permit any URL parameters.")
// 	}

// 	if path != "/" && path[len(path)-1] != '/' {
// 		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
// 		path += "/"
// 	}
// 	path += "*"

// 	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
// 		rctx := chi.RouteContext(r.Context())
// 		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
// 		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
// 		fs.ServeHTTP(w, r)
// 	})
// }

func setupFileServer(router *chi.Mux, root string) {
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
