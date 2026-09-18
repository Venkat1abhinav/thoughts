package main

import (
	"net/http"
	"strings"

	"github.com/owned_dragon/thoughts/gen/api"
	"github.com/owned_dragon/thoughts/internal/store"
)

func (app *application) GetFeed(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetFeedParams,
) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	if params.Limit != nil {
		fq.Limit = *params.Limit
	}

	if params.Offset != nil {
		fq.Offset = *params.Offset
	}

	if params.Search != nil {
		fq.Search = *params.Search
	}

	if params.Tags != nil {
		fq.Tags = strings.Split(*params.Tags, ",")
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	feed, err := app.store.Posts.GetUserFeed(r.Context(), 1, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
