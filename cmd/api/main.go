package main

import (
	"fmt"
	"github.com/rahuja23/GoBackendProject/internal/db"
	"github.com/rahuja23/GoBackendProject/internal/env"
	"github.com/rahuja23/GoBackendProject/internal/store"
	"log"
)

func main() {
	cfg := config{
		addr: env.GetString("API_ADDR", ":8000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social_networking_platform?sslmode=disable"),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 5),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	dbnew, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer dbnew.Close()
	fmt.Println("Database connection established")
	storage := store.NewStorage(dbnew)
	app := application{
		config: cfg,
		store:  storage,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))

}
