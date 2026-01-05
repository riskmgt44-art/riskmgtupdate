// middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"riskmgt/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.RespondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userName", claims.Name)
		ctx = context.WithValue(ctx, "userRole", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}