package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmonteiro/picpay-like/cmd/api"
	"github.com/jmonteiro/picpay-like/core/config"
	_ "github.com/lib/pq"
)

func main() {
	// Monta connection string do PostgreSQL
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBPort,
		config.Envs.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Database error:", err)
	}

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("DB: Successfully connected!")
}
