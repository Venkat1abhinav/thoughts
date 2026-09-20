package main

import (
	"net/http"

	"github.com/owned_dragon/thoughts/internal/store"
)

type RegisterUser struct {
	Username  string `json:"username" validate:"required,max=100"`
	FirstName string `json:"firstname" validate:"required,max=255"`
	LastName  string `json:"lastname" validate:"max=255"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) RegisterUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var payload RegisterUser
	ctx := r.Context()

	if err := readJSON(w, r, payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := store.User{
		Username:  payload.Username,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
	}

	if err := user.Password.
		Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Users.
		Create(ctx, &user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
