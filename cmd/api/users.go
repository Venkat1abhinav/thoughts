package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/owned_dragon/thoughts/internal/store"
)

type Follower struct {
	UserID int64 `json:"user_id"`
}

type UserKey string

const userkey UserKey = "user"

func (app *application) getUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := app.getUserContext(r)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		err = app.jsonResponse(w, http.StatusOK, user)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

	}
}

func (app *application) followUserHanlder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := app.getUserContext(r)
		if err != nil {
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
}

func (app *application) unfollowUserHanlder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := app.getUserContext(r)
		if err != nil {
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
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

			if err != nil {
				app.badRequestError(w, r, err)
				return
			}

			ctx := r.Context()

			user, err := app.store.Users.GetByID(ctx, id)

			if err != nil {
				switch {
				case errors.Is(err, store.ErrNotFound):
					app.notFoundError(w, r, err)
					return
				default:
					app.internalServerError(w, r, err)
					return
				}
			}

			ctx = context.WithValue(ctx, userkey, user)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)

		},
	)
}

func (app *application) getUserContext(r *http.Request) (*store.User, error) {
	user, ok := r.Context().Value(userkey).(*store.User)
	if !ok || user == nil {
		return nil, errors.New(
			"could not retrive user from user context",
		)
	}
	return user, nil

}
