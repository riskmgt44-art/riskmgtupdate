// handlers/audit_handler.go
package handlers

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"riskmgt/utils"
)

func GetAuditTrail(w http.ResponseWriter, r *http.Request) {
	targetID := r.URL.Query().Get("targetId")
	targetType := r.URL.Query().Get("type")

	query := bson.M{}
	if targetID != "" {
		objID, _ := primitive.ObjectIDFromHex(targetID)
		query["targetID"] = objID
	}
	if targetType != "" {
		query["targetType"] = targetType
	}

	cursor, err := auditCollection.Find(r.Context(), query, utils.PaginationOptions(r))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch audit trail")
		return
	}
	defer cursor.Close(r.Context())

	var audits []bson.M
	if err = cursor.All(r.Context(), &audits); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to decode audit records")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, audits)
}