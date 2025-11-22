package wallet

import (
	"fmt"
	"net/http"
	"strconv"
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
		r.Get("/user/{userID}", h.handleGetWalletsByUserID)
		r.Get("/{id}", h.handleGetWalletByID)
	})
}
//TODO: ONLY CARD PAYLOAD AND REGISTER REQUIRE JUST THE CARD NUMBER, BALANCE DEFAULT 0 
func (h *Handler) handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	userID := middelware.GetUserIDFromContext(r.Context())
	if userID == -1 {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	var payload struct {
		CardNumber string  `json:"card_number" validate:"required"`
		Balance    float64 `json:"balance"`
	}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	wallet := types.Wallet{
		UserID:     userID,
		CardNumber: payload.CardNumber,
		Balance:    payload.Balance,
	}

	err := h.walletService.CreateWallet(wallet)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "wallet created successfully"})
}

func (h *Handler) handleGetWalletByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	var id int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid wallet ID: %v", err))
		return
	}

	wallet, err := h.walletService.GetWalletByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("wallet not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, wallet)
}

func (h *Handler) handleGetWalletsByUserID(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "userID")
	var userID int
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID: %v", err))
		return
	}

	wallets, err := h.walletService.GetWalletsByUserID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("wallets not found for user"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, wallets)
}
