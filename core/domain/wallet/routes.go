package wallet

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	middelware "github.com/jmonteiro/picpay-like/core/middleware/auth"
	"github.com/jmonteiro/picpay-like/core/types"
	"github.com/jmonteiro/picpay-like/core/utils"
)

type Handler struct {
	userService   types.UserStore
	walletService *WalletService
}

func NewHandler(userService types.UserStore, walletService *WalletService) *Handler {
	return &Handler{
		userService:   userService,
		walletService: walletService,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/wallets", func(r chi.Router) {
		r.Post("/create", middelware.WithJWTAuth(h.handleCreateWallet, h.userService))
		//r.Get("/user/{userID}", h.handleGetWalletByUserID)
		//r.Get("/{id}", h.handleGetWalletByID)
	})
}

func (h *Handler) handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	// Pega o userID do contexto (usuário autenticado)
	userID := middelware.GetUserIDFromContext(r.Context())
	if userID == -1 {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	// Estrutura para receber apenas o balance inicial (opcional)
	var payload struct {
		Balance float64 `json:"balance"`
	}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Cria a carteira usando o userID do contexto
	wallet := types.Wallet{
		UserID:  userID,
		Balance: payload.Balance,
	}

	// Usa o service que já tem a lógica de verificação
	err := h.walletService.CreateWallet(wallet)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "wallet created successfully"})
}
