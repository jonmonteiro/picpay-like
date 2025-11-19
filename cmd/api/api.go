package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	userHandler "github.com/jmonteiro/picpay-like/core/handler/user"
	userRepo "github.com/jmonteiro/picpay-like/core/repository/user"
	userService "github.com/jmonteiro/picpay-like/core/service/user"
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

	r.Route("/api/v1", func(api chi.Router) {

		// ===== USER =====
		userStore := userRepo.NewStore(s.db)
		userSvc := userService.NewUserService(userStore)
		userHdlr := userHandler.NewHandler(userSvc)
		userHdlr.RegisterRoutes(api)

		// ===== WALLET =====
		// TODO: Implementar wallet handler

		// ===== TRANSACTION =====
		// TODO: Implementar transaction handler
	})

	// Servir arquivos estáticos (opcional)
	r.Handle("/*", http.FileServer(http.Dir("static")))

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, r)
}
