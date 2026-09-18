package main

import (
	"errors"
	"net/http"

	"github.com/owned_dragon/thoughts/internal/store"
)

type CommentCreate struct {
	Content string `json:"content" validate:"required"`
}

func (app *application) CreateComment(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}

		app.internalServerError(w, r, err)
		return
	}

	var commentCreate CommentCreate

	if err := readJSON(w, r, &commentCreate); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err = Validate.Struct(commentCreate); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	comment := &store.Comment{
		PostID:  post.ID,
		UserID:  post.UserID,
		Content: commentCreate.Content,
	}

	err = app.store.Comments.Create(r.Context(), comment)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(
		w,
		http.StatusCreated,
		comment,
	); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetComments(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	ctx := r.Context()

	comments, err := app.store.Comments.GetCommentsByPostID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}

		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusOK, comments)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
