// routes/approval_routes.go
package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"riskmgt/handlers"
	"riskmgt/middleware"
)

func RegisterApprovalRoutes(router *mux.Router) {
	approvals := router.PathPrefix("/approvals").Subrouter()

	approvals.HandleFunc("/pending", handlers.ListPendingApprovals).
		Methods("GET").
		Handler(middleware.HasRole("RiskManager")(http.HandlerFunc(handlers.ListPendingApprovals)))

	approvals.HandleFunc("/{id}/approve", handlers.ApproveItem).
		Methods("POST").
		Handler(middleware.HasRole("RiskManager")(http.HandlerFunc(handlers.ApproveItem)))

	approvals.HandleFunc("/{id}/reject", handlers.RejectItem).
		Methods("POST").
		Handler(middleware.HasRole("RiskManager")(http.HandlerFunc(handlers.RejectItem)))
}