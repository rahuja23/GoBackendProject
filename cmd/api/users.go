package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rahuja23/GoBackendProject/internal/store"
	"log"
	"net/http"
	"strconv"
)

type userKey string
type FollowUser struct {
	UserID int64 `json:"user_id"`
}

const userCtx userKey = "user"

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user *store.User
	ctx := r.Context()
	if err := readJSON(w, r, &user); err != nil {
		app.badrequestError(w, r, err)
		return

	}
	if err := Validate.Struct(user); err != nil {
		app.badrequestError(w, r, err)
		return
	}
	err := app.store.Users.Create(ctx, user)
	if err != nil {
		app.internalServerError(w, r, err)
	}
	if err := writeJSON(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	userId, err := strconv.Atoi(userID)
	if err != nil {
		log.Fatalf("Error Parsing UserID and converting to int: %s", err)
		return
	}
	ctx := r.Context()
	if err := app.store.Users.Delete(ctx, int64(userId)); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, "User Successfully Deleted"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	usrId, err := strconv.Atoi(userID)
	if err != nil {
		log.Fatalf("Error converting user id to int: %s\n", err)
		return
	}
	ctx := r.Context()
	getUser, err := app.store.Users.Get(ctx, int64(usrId))
	if err != nil {
		log.Fatalf("Error Fetching user from DB: %s\n", err)
		return
	}

	var payload *store.UpdateUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badrequestError(w, r, err)
		return

	}
	if payload.Username != "" || payload.Email != "" {
		if payload.Username != "" {
			getUser.Username = payload.Username
		} else {
			getUser.Email = payload.Email
		}
	} else {
		user_msg := fmt.Sprintf(" You can only update username or Email for User with ID %d ", usrId)
		app.jsonResponse(w, http.StatusNotFound, user_msg)
		return
	}
	user, err := app.store.Users.Update(ctx, getUser)
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) followHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}

	var follow FollowUser
	if err := readJSON(w, r, &follow); err != nil {
		app.badrequestError(w, r, err)
		return
	}
	ctx := r.Context()
	err := app.store.Followers.Follow(ctx, user.ID, follow.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) unfollowHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
	// revert back to auth.UserID
	var follow FollowUser
	if err := readJSON(w, r, &follow); err != nil {
		app.badrequestError(w, r, err)
		return
	}
	ctx := r.Context()
	err := app.store.Followers.Unfollow(ctx, user.ID, follow.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		usrId, err := strconv.Atoi(userID)
		if err != nil {
			log.Fatalf("Error converting user id to int: %s\n", err)
			return
		}
		ctx := r.Context()
		user, err := app.store.Users.Get(ctx, int64(usrId))
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notfoundError(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}

		}
		ctx = context.WithValue(r.Context(), userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
