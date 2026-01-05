// routes/audit_routes.go
package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"riskmgt/handlers"
	"riskmgt/middleware"
)

func RegisterAuditRoutes(router *mux.Router) {
	audit := router.PathPrefix("/audits").Subrouter()

	audit.HandleFunc("", handlers.GetAuditTrail).
		Methods("GET").
		Handler(middleware.HasRole("Executive")(http.HandlerFunc(handlers.GetAuditTrail)))
}