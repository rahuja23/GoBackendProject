package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rahuja23/GoBackendProject/internal/store"
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
	store  store.Storage
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	addr    string
	db      dbConfig
	env     string
	version string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Patch("/", app.updatePostHandler)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Route("/comment", func(r chi.Router) {
					r.Post("/", app.createCommentHandler)
					r.Delete("/{commentID}", app.deleteCommmentHandler)

					r.Put("/{commentID}", app.UpdateCommentHandler)
				})

			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Patch("/", app.updateUserHandler)
				r.Delete("/", app.deleteUserHandler)
				r.Put("/follow", app.followHandler)
				r.Put("/unfollow", app.unfollowHandler)
			})
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

	})

	return r
}

func (app *application) run(mux *chi.Mux) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("listening on %s", app.config.addr)

	return srv.ListenAndServe()
}
