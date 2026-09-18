package main

import (
	"errors"
	"net/http"

	"github.com/owned_dragon/thoughts/internal/store"
)

type Follower struct {
	UserID int64 `json:"user_id"`
}

func (app *application) GetUser(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	ctx := r.Context()
	user, err := app.store.Users.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}

		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) FollowUser(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	ctx := r.Context()
	user, err := app.store.Users.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}

		app.internalServerError(w, r, err)
		return
	}

	var follower Follower

	if err = readJSON(w, r, &follower); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.Followers.Follow(r.Context(), follower.UserID, user.ID)
	if err != nil {
		if errors.Is(err, store.ErrConfict) {
			app.confictError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusCreated, follower)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UnfollowUser(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	ctx := r.Context()
	user, err := app.store.Users.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}

		app.internalServerError(w, r, err)
		return
	}

	var follower Follower

	if err = readJSON(w, r, &follower); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.Followers.UnFollow(r.Context(), follower.UserID, user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusNoContent, follower)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
