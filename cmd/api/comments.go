package main

import (
	"net/http"

	"github.com/owned_dragon/thoughts/internal/store"
)

type CommentCreate struct {
	Content string `json:"content" validate:"required"`
}

func (app *application) createCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		post, err := app.getPostByContext(r)
		if err != nil {
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
		}

		comment := &store.Comment{
			PostID:  post.ID,
			UserID:  post.UserID,
			Content: commentCreate.Content,
		}

		err = app.store.Comments.Create(r.Context(), comment)
		if err != nil {
			app.internalServerError(w, r, err)
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
}

func (app *application) getCommentsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		post, err := app.getPostByContext(r)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		comments, err := app.store.Comments.
			GetCommentsByPostID(
				r.Context(),
				post.ID,
			)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		err = app.jsonResponse(w, http.StatusFound, comments)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
	}
}
