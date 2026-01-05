// utils/pagination.go
package utils

import (
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pagination struct {
	Page     int
	PageSize int
}

func ParsePagination(r *http.Request) Pagination {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return Pagination{Page: page, PageSize: pageSize}
}

func PaginationOptions(r *http.Request) *options.FindOptions {
	p := ParsePagination(r)
	skip := int64((p.Page - 1) * p.PageSize)
	limit := int64(p.PageSize)

	return options.Find().SetSkip(skip).SetLimit(limit)
}