package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rahuja23/GoBackendProject/internal/store"
	"log"
	"net/http"
	"strconv"
)

type CreateCommentPayload struct {
	Content string `json:"content"`
	UserId  int64  `json:"user_id"`
	PostId  int64  `json:"post_id"`
}
type UpdateCommentPayload struct {
	CommentId int64  `json:"id"`
	Content   string `json:"content"`
	UserId    int64  `json:"user_id"`
	PostId    int64  `json:"post_id"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	userId := 1
	id, err := strconv.ParseInt(idParam, 10, 64)
	ctx := r.Context()
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	var payload CreateCommentPayload
	payload.PostId = id
	payload.UserId = int64(userId)
	if err := readJSON(w, r, &payload); err != nil {
		app.badrequestError(w, r, err)
		return

	}
	if err := Validate.Struct(payload); err != nil {
		app.badrequestError(w, r, err)
		return
	}
	comment := &store.Comment{
		PostID:  payload.PostId,
		UserID:  payload.UserId,
		Content: payload.Content,
	}

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteCommmentHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	idParam2 := chi.URLParam(r, "commentID")

	post_id, err := strconv.ParseInt(idParam, 10, 64)
	comment_id, err := strconv.ParseInt(idParam2, 10, 64)
	ctx := r.Context()
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	comment := &store.CommentDelete{
		PostID: post_id,
		ID:     comment_id,
	}

	if err := app.store.Comments.Delete(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	idParam2 := chi.URLParam(r, "commentID")
	comment_id, err := strconv.ParseInt(idParam2, 10, 64)
	post_id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.internalServerError(w, r, err)
	}

	var payload UpdateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badrequestError(w, r, err)
		return

	}
	if err := Validate.Struct(payload); err != nil {
		app.badrequestError(w, r, err)
		return
	}
	fmt.Printf("Comment %d\n", comment_id)
	fmt.Printf("Post %d\n", post_id)
	fmt.Printf("Conent %s\n", payload.Content)
	fmt.Printf("User %d\n", payload.UserId)

	comment_mod := &store.Comment{
		ID:      comment_id,
		PostID:  post_id,
		UserID:  payload.UserId,
		Content: payload.Content,
	}

	if err := app.store.Comments.Update(r.Context(), comment_mod); err != nil {
		log.Printf("Error Updating the comment from Handlers side: %v", err)
		app.internalServerError(w, r, err)
	}

}
