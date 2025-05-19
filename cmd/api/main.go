package main

import (
	"github.com/rahuja23/GoBackendProject/internal/env"
	"github.com/rahuja23/GoBackendProject/internal/store"
	"log"
)

func main() {
	cfg := config{
		addr: env.GetString("API_ADDR", ":8080"),
	}
	storage := store.NewStorage(nil)
	app := application{
		config: cfg,
		store:  storage,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))

}
