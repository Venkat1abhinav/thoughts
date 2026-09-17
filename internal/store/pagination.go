package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"limit" validate:"gte=0`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate="max=5"`
	Search string   `json:"search" validate="max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	requestQuery := r.URL.Query()
	limit := requestQuery.Get("limit")

	if limit != "" {
		limit, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}
		fq.Limit = limit
	}

	offset := requestQuery.Get("offset")

	if offset != "" {
		offset, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = offset

	}

	sort := requestQuery.Get("sort")

	if sort != "" {
		fq.Sort = sort
	}

	tags := requestQuery.Get("tags")

	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := requestQuery.Get("search")

	if search != "" {
		fq.Search = search
	}

	since := requestQuery.Get("since")
	if since != "" {
		since, err := parseTime(since)
		if err != nil {
			return PaginatedFeedQuery{}, err
		}
		fq.Since = since
	}

	until := requestQuery.Get("until")
	if until != "" {
		until, err := parseTime(until)
		if err != nil {
			return PaginatedFeedQuery{}, err
		}
		fq.Until = until
	}

	return fq, nil
}

func parseTime(s string) (string, error) {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return "", err
	}

	return t.Format(time.DateTime), nil
}
