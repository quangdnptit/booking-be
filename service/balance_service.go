package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"booking-be/repo"
	"booking-be/repomodel"
)

var (
	ErrDepositAmountInvalid = errors.New("deposit amount must be positive")
	ErrUserNotFound         = errors.New("user not found")
)

// DepositResult is returned after a successful deposit.
type DepositResult struct {
	UserID       string  `json:"user_id"`
	DepositID    string  `json:"deposit_id"`
	Amount       float64 `json:"amount"`
	BalanceAfter float64 `json:"balance_after"`
	CreatedAt    string  `json:"created_at"`
}

// BalanceService handles user balance and deposits.
type BalanceService struct {
	users  repo.UserRepo
	ledger repo.LedgerRepo
}

func NewBalanceService(users repo.UserRepo, ledger repo.LedgerRepo) *BalanceService {
	return &BalanceService{users: users, ledger: ledger}
}

// Deposit adds amount to the user's balance and appends a balance record (same shape: amount, created_at).
func (s *BalanceService) Deposit(ctx context.Context, userID string, amount float64) (*DepositResult, error) {
	if amount <= 0 {
		return nil, ErrDepositAmountInvalid
	}
	rec, err := s.users.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if rec == nil {
		return nil, ErrUserNotFound
	}
	now := time.Now().UTC().Format(time.RFC3339)
	newBalance := rec.Amount + amount
	if err := s.users.AddBalance(ctx, userID, amount, now); err != nil {
		return nil, fmt.Errorf("update balance: %w", err)
	}
	depositID := uuid.New().String()
	ledgerRec := repomodel.BalanceRecord{
		UserID:       userID,
		DepositID:    depositID,
		Amount:       amount,
		BalanceAfter: newBalance,
		CreatedAt:    now,
	}
	if err := s.ledger.AddDeposit(ctx, ledgerRec); err != nil {
		return nil, fmt.Errorf("record deposit: %w", err)
	}
	return &DepositResult{
		UserID:       userID,
		DepositID:    depositID,
		Amount:       amount,
		BalanceAfter: newBalance,
		CreatedAt:    now,
	}, nil
}
