package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rahuja23/GoBackendProject/internal/store"
	"net/http"
	"strconv"
)

type CreateCommentPayload struct {
	Content string `json:"content"`
	UserId  int64  `json:"user_id"`
	PostId  int64  `json:"post_id"`
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
