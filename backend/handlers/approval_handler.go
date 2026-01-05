// handlers/approval_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"riskmgt/models"
	"riskmgt/services"
	"riskmgt/utils"
)

func ListPendingApprovals(w http.ResponseWriter, r *http.Request) {
	userRole := r.Context().Value("userRole").(string)
	if userRole != "RiskManager" && userRole != "Executive" {
		utils.RespondWithError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	query := bson.M{"status": "Pending"}
	cursor, err := approvalCollection.Find(r.Context(), query)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch approvals")
		return
	}
	defer cursor.Close(r.Context())

	var approvals []models.Approval
	if err = cursor.All(r.Context(), &approvals); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to decode approvals")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, approvals)
}

func ApproveItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	userID := r.Context().Value("userID").(primitive.ObjectID)
	userRole := r.Context().Value("userRole").(string)

	if err := services.ProcessApproval(r.Context(), id, "Approved", userID, userRole); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Item approved"})
}

func RejectItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	var payload struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Reason required")
		return
	}

	userID := r.Context().Value("userID").(primitive.ObjectID)
	userRole := r.Context().Value("userRole").(string)

	if err := services.ProcessApproval(r.Context(), id, "Rejected", userID, userRole, payload.Reason); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Item rejected"})
}