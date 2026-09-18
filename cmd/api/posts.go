package main

import (
	"errors"
	"net/http"

	"github.com/owned_dragon/thoughts/internal/store"
)

type PostCreate struct {
	Content string   `json:"content" validate:"required,max=100"`
	Title   string   `json:"title" validate:"required,max=1000"`
	UserID  int64    `json:"user_id"`
	Tags    []string `json:"tags"`
}

type PostUpdate struct {
	Content *string `json:"content" validate:"omitempty,max=100"`
	Title   *string `json:"title" validate:"omitempty,max=1000"`
}

func (app *application) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post PostCreate

	if err := readJSON(w, r, &post); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(post); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	createdPost := &store.Post{
		Content: post.Content,
		Title:   post.Title,
		UserID:  post.UserID,
		Tags:    post.Tags,
	}

	if err := app.store.Posts.Create(ctx, createdPost); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(
		w,
		http.StatusCreated,
		createdPost,
	); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetPost(
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
	comments, err := app.store.Comments.GetCommentsByPostID(
		r.Context(),
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

	post.Comments = comments

	if err := app.jsonResponse(
		w,
		http.StatusOK,
		post,
	); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UpdatePost(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	post, err := app.store.Posts.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}

		app.internalServerError(w, r, err)
		return
	}

	var input PostUpdate

	if err := readJSON(w, r, &input); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(input); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if input.Title == nil && input.Content == nil {
		app.badRequestError(
			w,
			r,
			errors.New("at least one field must be provided"),
		)
		return
	}

	if input.Title != nil {
		post.Title = *input.Title
	}

	if input.Content != nil {
		post.Content = *input.Content
	}

	updatedPost, err := app.store.Posts.UpdateByID(
		r.Context(),
		post,
	)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
			return
		case errors.Is(err, store.ErrVersionConflict):
			app.confictError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	comments, err := app.store.Comments.GetCommentsByPostID(
		r.Context(),
		post.ID,
	)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	updatedPost.Comments = comments

	if err := app.jsonResponse(
		w,
		http.StatusOK,
		updatedPost,
	); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeletePost(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	ctx := r.Context()
	err := app.store.Posts.DeleteByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
