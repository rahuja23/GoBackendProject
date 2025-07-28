package main

import (
	"github.com/rahuja23/GoBackendProject/internal/db"
	"github.com/rahuja23/GoBackendProject/internal/env"
	"github.com/rahuja23/GoBackendProject/internal/store"
	"log"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social_networking_platform?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := store.NewStorage(conn)
	db.Seed(store)
}
