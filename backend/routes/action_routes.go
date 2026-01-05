// routes/action_routes.go
package routes

import (
	"github.com/gorilla/mux"
	"riskmgt/handlers"
)

func RegisterActionRoutes(router *mux.Router) {
	actions := router.PathPrefix("/actions").Subrouter()

	actions.HandleFunc("", handlers.CreateAction).Methods("POST")
	actions.HandleFunc("", handlers.ListActions).Methods("GET")
	actions.HandleFunc("/{id}", handlers.GetAction).Methods("GET")
	actions.HandleFunc("/{id}", handlers.UpdateAction).Methods("PUT")
	actions.HandleFunc("/{id}/submit", handlers.SubmitActionForApproval).Methods("POST")
}