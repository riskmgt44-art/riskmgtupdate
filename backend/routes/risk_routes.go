// routes/risk_routes.go
package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"riskmgt/handlers"
	"riskmgt/middleware"
)

func RegisterRiskRoutes(router *mux.Router) {
	risks := router.PathPrefix("/risks").Subrouter()

	risks.HandleFunc("", handlers.CreateRisk).Methods("POST")
	risks.HandleFunc("", handlers.ListRisks).Methods("GET")
	risks.HandleFunc("/{id}", handlers.GetRisk).Methods("GET")
	risks.HandleFunc("/{id}", handlers.UpdateRisk).Methods("PUT")
	risks.HandleFunc("/{id}/submit", handlers.SubmitRiskForApproval).Methods("POST")

	// RAM update restricted to RiskManager+
	risks.HandleFunc("/{id}/ram", handlers.UpdateRiskRAM).
		Methods("POST").
		Handler(middleware.HasRole("RiskManager")(http.HandlerFunc(handlers.UpdateRiskRAM)))
}