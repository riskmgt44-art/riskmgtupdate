// handlers/action_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"riskmgt/models"
	"riskmgt/services"
	"riskmgt/utils"
)

func CreateAction(w http.ResponseWriter, r *http.Request) {
	var action models.Action
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	action.ID = primitive.NewObjectID()
	action.Status = "Draft"
	action.CreatedAt = time.Now()
	action.UpdatedAt = time.Now()
	action.OwnerID = r.Context().Value("userID").(primitive.ObjectID)
	action.OwnerName = r.Context().Value("userName").(string)

	if err := services.ValidateActionCreation(&action); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := actionCollection.InsertOne(r.Context(), action)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create action")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, action)
}

func UpdateAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	var updateData models.ActionUpdate
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var action models.Action
	err := actionCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&action)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Action not found")
		return
	}

	if action.Status != "Draft" {
		utils.RespondWithError(w, http.StatusForbidden, "Only draft actions can be edited")
		return
	}

	userID := r.Context().Value("userID").(primitive.ObjectID)
	if action.OwnerID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "You can only edit your own actions")
		return
	}

	update := bson.M{"$set": bson.M{
		"title":       updateData.Title,
		"description": updateData.Description,
		"dueDate":     updateData.DueDate,
		"assigneeID":  updateData.AssigneeID,
		"updatedAt":   time.Now(),
	}}

	_, err = actionCollection.UpdateOne(r.Context(), bson.M{"_id": id}, update)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update action")
		return
	}

	actionCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&action)
	utils.RespondWithJSON(w, http.StatusOK, action)
}

func SubmitActionForApproval(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	var action models.Action
	err := actionCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&action)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Action not found")
		return
	}

	if action.Status != "Draft" {
		utils.RespondWithError(w, http.StatusBadRequest, "Only draft actions can be submitted")
		return
	}

	userID := r.Context().Value("userID").(primitive.ObjectID)
	if action.OwnerID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "You can only submit your own actions")
		return
	}

	if err := services.SubmitActionForApproval(r.Context(), &action); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	actionCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&action)
	utils.RespondWithJSON(w, http.StatusOK, action)
}

func GetAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	var action models.Action
	err := actionCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&action)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Action not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, action)
}

func ListActions(w http.ResponseWriter, r *http.Request) {
	filters := utils.ParseQueryFilters(r)
	page := utils.ParsePagination(r)

	userRole := r.Context().Value("userRole").(string)
	userID := r.Context().Value("userID").(primitive.ObjectID)

	query := bson.M{}
	switch userRole {
	case "Analyst":
		query["ownerID"] = userID
	case "Viewer":
		query["status"] = "Approved"
	}

	actions, total, err := services.PaginateActions(r.Context(), query, filters, page)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch actions")
		return
	}

	totalPages := int((total + int64(page.PageSize) - 1) / int64(page.PageSize))

	response := map[string]any{
		"data":       actions,
		"total":      total,
		"page":       page.Page,
		"pageSize":   page.PageSize,
		"totalPages": totalPages,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}