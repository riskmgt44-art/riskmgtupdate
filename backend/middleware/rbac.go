// middleware/rbac.go
package middleware

import (
	"net/http"

	"riskmgt/utils"
)

var roleHierarchy = map[string]int{
	"Viewer":       1,
	"Analyst":      2,
	"RiskManager":  3,
	"Executive":    4,
	"Admin":        5,
}

func HasRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("userRole").(string)
			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "Role not found in context")
				return
			}

			userLevel := roleHierarchy[userRole]
			requiredLevel := roleHierarchy[requiredRole]

			if userLevel < requiredLevel {
				utils.RespondWithError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Public routes (login, forgot password) skip auth
func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}
		// If token present, try to authenticate (best effort)
		AuthMiddleware(next).ServeHTTP(w, r)
	})
}