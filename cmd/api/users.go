package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/owned_dragon/thoughts/internal/store"
)

type Userkey string

const userKey Userkey = "user"

func (app *application) getUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := app.GetUserContext(r)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		err = app.jsonResponse(w, http.StatusFound, user)

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

			ctx = context.WithValue(ctx, userKey, user)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)

		},
	)
}

func (app *application) GetUserContext(r *http.Request) (*store.User, error) {
	user, ok := r.Context().Value(userKey).(*store.User)
	if !ok || user == nil {
		return nil, errors.New(
			"could not retrive user from user context",
		)
	}

	return user, nil

}
