// Package http contains the implementation of the HTTP server.
package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akazantzidis/gwi-ass/internal/app"
	"github.com/akazantzidis/gwi-ass/internal/infra/http/favourite"
	"github.com/gorilla/mux"
)

// Server Represents the http server running for this service
type Server struct {
	appServicesF app.Services
	router       *mux.Router
}

// NewServer HTTP Server constructor
func NewServer(appServicesF app.Services) *Server {
	httpServer := &Server{appServicesF: appServicesF}
	httpServer.router = mux.NewRouter()
	httpServer.AddFavoriteHTTPRoutes()
	http.Handle("/", httpServer.router)

	return httpServer
}

func (httpServer *Server) AddFavoriteHTTPRoutes() {
	const base = "/users/{userID}/favorites"

	h := favourite.NewHandlerFF(httpServer.appServicesF.FavoriteServices)

	// Queries
	httpServer.router.HandleFunc(
		base,
		h.GetAllF,
	).Methods("GET")

	httpServer.router.HandleFunc(base+"/{favoriteId}",
		h.GetByIDF,
	).Methods("GET")

	// Commands
	httpServer.router.HandleFunc(base,
		h.Create,
	).Methods("POST")

	httpServer.router.HandleFunc(base+"/{favoriteId}",
		h.Patch).Methods("PATCH")

	httpServer.router.HandleFunc(base+"/{favoriteId}",
		h.Update,
	).Methods("PUT")

	httpServer.router.HandleFunc(base+"/{favoriteId}",
		h.Delete,
	).Methods("DELETE")
}

// ListenAndServe Starts listening for requests
func (httpServer *Server) ListenAndServe(port string) {
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
