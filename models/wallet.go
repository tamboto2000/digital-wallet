package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Wallet represents a user's digital wallet
type Wallet struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Transaction represents a wallet transaction
type Transaction struct {
	ID            int64     `json:"id"`
	WalletID      int64     `json:"wallet_id"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"` // "withdraw" or "deposit"
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
}

var (
	ErrWalletNotFound     = errors.New("wallet not found")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInvalidUserID      = errors.New("invalid user ID")
)

// GetWalletByUserID retrieves a wallet by user ID
func GetWalletByUserID(db *sql.DB, userID int64) (*Wallet, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	wallet := &Wallet{}
	query := `
		SELECT id, user_id, balance, created_at, updated_at 
		FROM wallets 
		WHERE user_id = $1
	`
	
	err := db.QueryRow(query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrWalletNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching wallet: %w", err)
	}
	
	return wallet, nil
}

// Withdraw performs a withdrawal from the wallet with transaction support
func Withdraw(db *sql.DB, userID int64, amount float64) (*Wallet, *Transaction, error) {
	// Validate amount
	if amount <= 0 {
		return nil, nil, ErrInvalidAmount
	}

	// Start transaction to ensure atomicity
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if committed

	// Lock the wallet row for update to prevent race conditions
	wallet := &Wallet{}
	query := `
		SELECT id, user_id, balance, created_at, updated_at 
		FROM wallets 
		WHERE user_id = $1 
		FOR UPDATE
	`
	
	err = tx.QueryRow(query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil, ErrWalletNotFound
	}
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching wallet: %w", err)
	}

	// Check if sufficient funds
	if wallet.Balance < amount {
		return nil, nil, ErrInsufficientFunds
	}

	balanceBefore := wallet.Balance
	balanceAfter := wallet.Balance - amount

	// Update wallet balance
	updateQuery := `
		UPDATE wallets 
		SET balance = $1, updated_at = $2 
		WHERE id = $3
	`
	_, err = tx.Exec(updateQuery, balanceAfter, time.Now(), wallet.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("error updating wallet: %w", err)
	}

	// Record transaction
	transaction := &Transaction{
		WalletID:      wallet.ID,
		Amount:        amount,
		Type:          "withdraw",
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		CreatedAt:     time.Now(),
	}

	insertQuery := `
		INSERT INTO transactions (wallet_id, amount, type, balance_before, balance_after, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err = tx.QueryRow(
		insertQuery,
		transaction.WalletID,
		transaction.Amount,
		transaction.Type,
		transaction.BalanceBefore,
		transaction.BalanceAfter,
		transaction.CreatedAt,
	).Scan(&transaction.ID)
	
	if err != nil {
		return nil, nil, fmt.Errorf("error recording transaction: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("error committing transaction: %w", err)
	}

	// Update wallet object
	wallet.Balance = balanceAfter
	wallet.UpdatedAt = time.Now()

	return wallet, transaction, nil
}
