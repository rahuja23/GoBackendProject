package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rahuja23/GoBackendProject/internal/store"
	"net/http"
	"strconv"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource conflict")
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type UpdatePost struct {
	Title   *string `json:"title" validate:"omitempty,max=700"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
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
		Version: 0,
	}

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)

	}
}
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var payload UpdatePost

	if err := readJSON(w, r, &payload); err != nil {
		app.badrequestError(w, r, err)
		return

	}

	if err := Validate.Struct(payload); err != nil {
		app.badrequestError(w, r, err)
	}
	//fmt.Println(*payload.Title)
	//fmt.Println(*payload.Content)

	ctx := r.Context()
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(ctx, post); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
	return
}
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	comments, err := app.store.Comments.Get(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments
	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)

	}
}
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	post_id, err := strconv.ParseInt(idParam, 10, 64)
	ctx := r.Context()
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Posts.Delete(ctx, post_id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notfoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
			return

		}
		w.WriteHeader(http.StatusNoContent)
	}

}
func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		ctx := r.Context()
		post, err := app.store.Posts.Get(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notfoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
		}
		ctx = context.WithValue(r.Context(), postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
