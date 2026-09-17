package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"limit" validate:"gte=0`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {

	limit := r.URL.Query().Get("limit")

	if limit != "" {
		limit, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}
		fq.Limit = limit
	}

	offset := r.URL.Query().Get("offset")

	if offset != "" {
		offset, err := strconv.Atoi(offset)

		if err != nil {
			return fq, nil
		}
		fq.Offset = offset

	}

	sort := r.URL.Query().Get("sort")

	if sort != "" {
		fq.Sort = sort
	}

	return fq, nil
}
