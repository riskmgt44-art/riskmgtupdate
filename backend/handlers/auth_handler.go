// handlers/auth_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"riskmgt/models"
	"riskmgt/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	var user models.User
	err := userCollection.FindOne(r.Context(), bson.M{"email": creds.Email}).Decode(&user)
	if err != nil || !utils.CheckPasswordHash(creds.Password, user.PasswordHash) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Role, user.Name)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := map[string]any{
		"token": token,
		"user": map[string]any{
			"id":   user.ID.Hex(),
			"name": user.Name,
			"role": user.Role,
		},
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Reset link sent"})
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Password reset successful"})
}