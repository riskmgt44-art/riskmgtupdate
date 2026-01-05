// handlers/user_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"riskmgt/models"
	"riskmgt/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	user.ID = primitive.NewObjectID()
	hash, _ := utils.HashPassword(user.PasswordHash)
	user.PasswordHash = hash

	_, err := userCollection.InsertOne(r.Context(), user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	user.PasswordHash = ""
	utils.RespondWithJSON(w, http.StatusCreated, user)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	cursor, err := userCollection.Find(r.Context(), bson.M{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	defer cursor.Close(r.Context())

	var users []models.User
	if err = cursor.All(r.Context(), &users); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Decode error")
		return
	}

	for i := range users {
		users[i].PasswordHash = ""
	}

	utils.RespondWithJSON(w, http.StatusOK, users)
}