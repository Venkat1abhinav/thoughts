package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/owned_dragon/thoughts/internal/store"
)

type PostCreate struct {
	Content string   `json:"content" validate:"required,max=100"`
	Title   string   `json:"title" validate:"required,max=1000"`
	UserID  int64    `json:"user_id"`
	Tags    []string `json:"tags"`
}

type contextKey string

const postKey contextKey = "post"

type PostUpdate struct {
	Content *string `json:"content" validate:"omitempty,max=100"`
	Title   *string `json:"title" validate:"omitempty,max=1000"`
}

func (app *application) createPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func (app *application) getPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		post, err := app.getPostByContext(r)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		comments, err := app.store.Comments.GetPostByID(
			r.Context(),
			post.ID,
		)
		if err != nil {
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
}

func (app *application) deletePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(
			chi.URLParam(r, "id"),
			10,
			64,
		)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()

		if err := app.store.Posts.DeleteByID(ctx, id); err != nil {
			if errors.Is(err, store.ErrNotFound) {
				app.notFoundError(w, r, err)
				return
			}

			app.internalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (app *application) updatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		post, err := app.getPostByContext(r)
		if err != nil {
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
		post.Content = *input.Content

		updatedPost, err := app.store.Posts.UpdateByID(
			r.Context(),
			post,
		)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				app.notFoundError(w, r, err)
				return
			}

			app.internalServerError(w, r, err)
			return
		}

		comments, err := app.store.Comments.GetPostByID(
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
}

func (app *application) postsContextMiddleware(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseInt(
				chi.URLParam(r, "id"),
				10,
				64,
			)
			if err != nil {
				app.badRequestError(w, r, err)
				return
			}

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

			ctx = context.WithValue(ctx, postKey, post)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		},
	)
}

func (app *application) getPostByContext(
	r *http.Request,
) (*store.Post, error) {
	post, ok := r.Context().Value(postKey).(*store.Post)
	if !ok || post == nil {
		return nil, errors.New(
			"could not retrieve post from request context",
		)
	}

	return post, nil
}
