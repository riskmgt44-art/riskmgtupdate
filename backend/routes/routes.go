// routes/routes.go
package routes

import (
	"github.com/gorilla/mux"
	"riskmgt/handlers"
	"riskmgt/middleware"
)

func RegisterRoutes(router *mux.Router) {
	// Public routes - no auth required
	public := router.PathPrefix("/api/auth").Subrouter()
	public.Use(middleware.OptionalAuth)
	public.HandleFunc("/login", handlers.Login).Methods("POST")
	public.HandleFunc("/forgot", handlers.ForgotPassword).Methods("POST")
	public.HandleFunc("/reset", handlers.ResetPassword).Methods("POST")

	// Protected routes - require JWT
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// Register all protected routes
	RegisterRiskRoutes(api)
	RegisterActionRoutes(api)
	RegisterApprovalRoutes(api)
	RegisterAuditRoutes(api)
	RegisterAdminRoutes(api)
	// Dashboard uses existing endpoints with role filtering in handlers
}