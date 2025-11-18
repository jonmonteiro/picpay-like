package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/jmonteiro/picpay-like/core/domain/user"
	"github.com/jmonteiro/picpay-like/core/domain/wallet"
	"github.com/jmonteiro/picpay-like/core/domain/transaction"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Prefixo da API
	r.Route("/api/v1", func(api chi.Router) {

		// ===== USER =====
		userStore := user.NewStore(s.db)
		userHandler := user.NewHandler(userStore)
		userHandler.RegisterRoutes(api)

		// ===== WALLET =====
		walletStore := wallet.NewStore(s.db)
		walletHandler := wallet.NewHandler(walletStore, userStore)
		walletHandler.RegisterRoutes(api)

		// ===== TRANSACTION =====
		transactionStore := transaction.NewStore(s.db)
		transactionHandler := transaction.NewHandler(transactionStore, walletStore, userStore)
		transactionHandler.RegisterRoutes(api)
	})

	// Servir arquivos estáticos (opcional)
	r.Handle("/*", http.FileServer(http.Dir("static")))

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, r)
}
