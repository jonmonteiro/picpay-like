package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	domain "github.com/jmonteiro/picpay-like/core/domain/user"
	service "github.com/jmonteiro/picpay-like/core/service/user"
	"github.com/jmonteiro/picpay-like/core/utils"
)

type Handler struct {
	userService *service.UserService
}

func NewHandler(userService *service.UserService) *Handler {
	return &Handler{userService: userService}
}

// RegisterRoutes registra as rotas de usuários
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/register", h.handleRegister) // POST /api/v1/users
		//r.Post("/login", h.handleCreateUser)
		//r.Get("/", h.handleListUsers)   // GET /api/v1/users
		//r.Get("/{id}", h.handleGetUser) // GET /api/v1/users/:id
	})
}

// handleRegister cria um novo usuário
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user domain.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.userService.RegisterUser(user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user created successfully"})
}
