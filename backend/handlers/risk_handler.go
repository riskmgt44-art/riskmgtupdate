// handlers/risk_handler.go
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

func CreateRisk(w http.ResponseWriter, r *http.Request) {
	var risk models.Risk
	if err := json.NewDecoder(r.Body).Decode(&risk); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	risk.ID = primitive.NewObjectID()
	risk.Status = "Draft"
	risk.CreatedAt = time.Now()
	risk.UpdatedAt = time.Now()
	risk.OwnerID = r.Context().Value("userID").(primitive.ObjectID)
	risk.OwnerName = r.Context().Value("userName").(string)

	if err := services.ValidateRiskCreation(&risk); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := riskCollection.InsertOne(r.Context(), risk)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create risk")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, risk)
}

func UpdateRisk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	var updateData models.RiskUpdate
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var risk models.Risk
	err := riskCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&risk)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Risk not found")
		return
	}

	if risk.Status != "Draft" {
		utils.RespondWithError(w, http.StatusForbidden, "Only draft risks can be edited")
		return
	}

	userID := r.Context().Value("userID").(primitive.ObjectID)
	if risk.OwnerID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "You can only edit your own risks")
		return
	}

	update := bson.M{"$set": bson.M{
		"title":              updateData.Title,
		"description":        updateData.Description,
		"category":           updateData.Category,
		"likelihood":         updateData.Likelihood,
		"impact":             updateData.Impact,
		"residualLikelihood": updateData.ResidualLikelihood,
		"residualImpact":     updateData.ResidualImpact,
		"updatedAt":          time.Now(),
	}}

	_, err = riskCollection.UpdateOne(r.Context(), bson.M{"_id": id}, update)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update risk")
		return
	}

	riskCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&risk)
	utils.RespondWithJSON(w, http.StatusOK, risk)
}

func SubmitRiskForApproval(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	var risk models.Risk
	err := riskCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&risk)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Risk not found")
		return
	}

	if risk.Status != "Draft" {
		utils.RespondWithError(w, http.StatusBadRequest, "Only draft risks can be submitted")
		return
	}

	userID := r.Context().Value("userID").(primitive.ObjectID)
	if risk.OwnerID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "You can only submit your own risks")
		return
	}

	if err := services.SubmitRiskForApproval(r.Context(), &risk); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	riskCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&risk)
	utils.RespondWithJSON(w, http.StatusOK, risk)
}

func GetRisk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	var risk models.Risk
	err := riskCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&risk)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Risk not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, risk)
}

func ListRisks(w http.ResponseWriter, r *http.Request) {
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

	risks, total, err := services.PaginateRisks(r.Context(), query, filters, page)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch risks")
		return
	}

	totalPages := int((total + int64(page.PageSize) - 1) / int64(page.PageSize))

	response := map[string]any{
		"data":       risks,
		"total":      total,
		"page":       page.Page,
		"pageSize":   page.PageSize,
		"totalPages": totalPages,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

func UpdateRiskRAM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	var ramUpdate models.RAMUpdate
	if err := json.NewDecoder(r.Body).Decode(&ramUpdate); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	userID := r.Context().Value("userID").(primitive.ObjectID)
	if err := services.UpdateRiskRAM(r.Context(), id, &ramUpdate, userID); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var risk models.Risk
	riskCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&risk)
	utils.RespondWithJSON(w, http.StatusOK, risk)
}