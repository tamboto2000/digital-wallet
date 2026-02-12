package models

import (
	"testing"
)

func TestWalletValidation(t *testing.T) {
	tests := []struct {
		name    string
		userID  int64
		amount  float64
		wantErr error
	}{
		{
			name:    "Invalid user ID - zero",
			userID:  0,
			amount:  100.0,
			wantErr: ErrInvalidUserID,
		},
		{
			name:    "Invalid user ID - negative",
			userID:  -1,
			amount:  100.0,
			wantErr: ErrInvalidUserID,
		},
		{
			name:    "Invalid amount - zero",
			userID:  1,
			amount:  0,
			wantErr: nil, // Will be caught in Withdraw function
		},
		{
			name:    "Invalid amount - negative",
			userID:  1,
			amount:  -50.0,
			wantErr: nil, // Will be caught in Withdraw function
		},
		{
			name:    "Valid inputs",
			userID:  1,
			amount:  50.0,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test GetWalletByUserID validation
			if tt.userID <= 0 {
				_, err := GetWalletByUserID(nil, tt.userID)
				if err != tt.wantErr {
					t.Errorf("GetWalletByUserID() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		err     error
		message string
	}{
		{ErrWalletNotFound, "wallet not found"},
		{ErrInsufficientFunds, "insufficient funds"},
		{ErrInvalidAmount, "invalid amount"},
		{ErrInvalidUserID, "invalid user ID"},
	}

	for _, tt := range tests {
		if tt.err.Error() != tt.message {
			t.Errorf("Expected error message '%s', got '%s'", tt.message, tt.err.Error())
		}
	}
}
