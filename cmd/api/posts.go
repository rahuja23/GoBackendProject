package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rahuja23/GoBackendProject/internal/store"
	"net/http"
	"strconv"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=200"`
	Content string   `json:"content" validate:"required"`
	Tags    []string `json:"tags"`
}

type EditPost struct {
	Content string `json:"content" validate:"required"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	userId := 1
	ctx := r.Context()

	if err := readJSON(w, r, &payload); err != nil {
		app.badrequestError(w, r, err)
		return

	}
	if err := Validate.Struct(payload); err != nil {
		app.badrequestError(w, r, err)
		return
	}
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  int64(userId),
		Tags:    payload.Tags,
	}

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)

	}

}

func (app *application) editPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	var payload EditPost
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := readJSON(w, r, &payload); err != nil {
		app.badrequestError(w, r, err)
		return

	}
	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		app.notfoundError(w, r, err)
	}
	post.Content = payload.Content
	post, err = app.store.Posts.UpdateByID(ctx, post)
	fmt.Println("Post Content Now: ", post.Content)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notfoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}
	comments, err := app.store.Comments.GetCommentsByPostId(ctx, id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)

	}
}
