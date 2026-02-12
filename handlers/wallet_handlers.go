package handlers

import (
	"database/sql"
	"digital-wallet/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Response structures
type ErrorResponse struct {
	Error string `json:"error"`
}

type BalanceResponse struct {
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount"`
}

type WithdrawResponse struct {
	UserID        int64   `json:"user_id"`
	Amount        float64 `json:"amount"`
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
	TransactionID int64   `json:"transaction_id"`
}

// GetBalance handles GET /api/balance/{user_id}
func GetBalance(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract user_id from URL
		vars := mux.Vars(r)
		userIDStr := vars["user_id"]
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid user ID"})
			return
		}

		// Get wallet
		wallet, err := models.GetWalletByUserID(db, userID)
		if err != nil {
			if errors.Is(err, models.ErrWalletNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Wallet not found"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			return
		}

		// Return balance
		response := BalanceResponse{
			UserID:  wallet.UserID,
			Balance: wallet.Balance,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// Withdraw handles POST /api/withdraw/{user_id}
func Withdraw(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract user_id from URL
		vars := mux.Vars(r)
		userIDStr := vars["user_id"]
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid user ID"})
			return
		}

		// Parse request body
		var req WithdrawRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
			return
		}

		// Perform withdrawal
		wallet, transaction, err := models.Withdraw(db, userID, req.Amount)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrWalletNotFound):
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Wallet not found"})
			case errors.Is(err, models.ErrInsufficientFunds):
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Insufficient funds"})
			case errors.Is(err, models.ErrInvalidAmount):
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Amount must be greater than 0"})
			default:
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			}
			return
		}

		// Return success response
		response := WithdrawResponse{
			UserID:        wallet.UserID,
			Amount:        req.Amount,
			BalanceBefore: transaction.BalanceBefore,
			BalanceAfter:  wallet.Balance,
			TransactionID: transaction.ID,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
