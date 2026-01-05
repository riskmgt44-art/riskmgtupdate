// routes/admin_routes.go
package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"riskmgt/handlers"
	"riskmgt/middleware"
)

func RegisterAdminRoutes(router *mux.Router) {
	admin := router.PathPrefix("/admin").Subrouter()

	admin.HandleFunc("/users", handlers.CreateUser).
		Methods("POST").
		Handler(middleware.HasRole("Executive")(http.HandlerFunc(handlers.CreateUser)))

	admin.HandleFunc("/users", handlers.ListUsers).
		Methods("GET").
		Handler(middleware.HasRole("Executive")(http.HandlerFunc(handlers.ListUsers)))
}