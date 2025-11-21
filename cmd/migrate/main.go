package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/jmonteiro/picpay-like/core/config"
)

func main() {
	// Monta connection string do PostgreSQL
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBHost,
		config.Envs.DBPort,
		config.Envs.DBName,
	)

	// Conecta ao banco de dados
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verifica se a conexão está funcionando
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Cria a instância do driver de migração do PostgreSQL
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Cria a instância de migração
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Mostra a versão atual
	v, d, _ := m.Version()
	log.Printf("Version: %d, dirty: %v", v, d)

	// Pega o comando da linha de comando
	cmd := os.Args[len(os.Args)-1]

	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migration UP executada com sucesso!")
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migration DOWN executada com sucesso!")
	}

	if cmd == "force" {
		if len(os.Args) < 3 {
			log.Fatal("Uso: go run cmd/migrate/main.go force <version>")
		}
		version := os.Args[2]
		var forceVersion int
		fmt.Sscanf(version, "%d", &forceVersion)
		if err := m.Force(forceVersion); err != nil {
			log.Fatal(err)
		}
		log.Printf("Forçada versão %d com sucesso!", forceVersion)
	}
}
