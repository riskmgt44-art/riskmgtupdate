// utils/query.go
package utils

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func ParseQueryFilters(r *http.Request) bson.M {
	filters := bson.M{}
	query := r.URL.Query()

	if status := query.Get("status"); status != "" {
		filters["status"] = status
	}
	if category := query.Get("category"); category != "" {
		filters["category"] = category
	}
	if search := query.Get("search"); search != "" {
		regex := bson.M{"$regex": search, "$options": "i"}
		filters["$or"] = []bson.M{
			{"title": regex},
			{"description": regex},
		}
	}

	return filters
}