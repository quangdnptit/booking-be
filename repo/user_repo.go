package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	dynamo "github.com/guregu/dynamo/v2"

	"booking-be/repomodel"
)

const (
	TableUsers = "users"
	IndexEmail = "email-index"
)

// Single table: user profile (pk=USER#id, sk=PROFILE) and deposit rows (pk=USER#id, sk=DEPOSIT#id).
// GSI "email-index": partition key = email, sort key = sk.

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*repomodel.UserRecord, error)
	GetByUserID(ctx context.Context, userID string) (*repomodel.UserRecord, error)
	Create(ctx context.Context, rec repomodel.UserRecord) error
	AddBalance(ctx context.Context, userID string, delta float64, updatedAt string) error
}

type DynamoUserRepo struct {
	table dynamo.Table
}

func NewDynamoUserRepo(db *dynamo.DB) *DynamoUserRepo {
	return &DynamoUserRepo{table: db.Table(TableUsers)}
}

func normalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func userPK(userID string) string { return repomodel.PKPrefixUser + userID }

// GetByEmail queries GSI "email-index" (PK=email, SK=sk). Returns the profile row.
func (r *DynamoUserRepo) GetByEmail(ctx context.Context, email string) (*repomodel.UserRecord, error) {
	key := normalizeEmail(email)
	if key == "" {
		return nil, nil
	}
	var rec repomodel.UserRecord
	err := r.table.Get("email", key).Index(IndexEmail).One(ctx, &rec)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &rec, nil
}

// GetByUserID gets the profile row by pk
func (r *DynamoUserRepo) GetByUserID(ctx context.Context, userID string) (*repomodel.UserRecord, error) {
	if userID == "" {
		return nil, nil
	}
	var rec repomodel.UserRecord
	err := r.table.Get("pk", userPK(userID)).One(ctx, &rec)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &rec, nil
}

// Create puts the user profile. Email uniqueness enforced via GetByEmail (same table, GSI).
func (r *DynamoUserRepo) Create(ctx context.Context, rec repomodel.UserRecord) error {
	rec.Email = normalizeEmail(rec.Email)
	if rec.Email == "" {
		return fmt.Errorf("email required")
	}
	existing, err := r.GetByEmail(ctx, rec.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("email already registered")
	}
	rec.Pk = userPK(rec.UserID)
	rec.Sk = repomodel.SKProfile
	return r.table.Put(rec).Run(ctx)
}

// AddBalance adds depositAmount to the user's amount (by user_id).
func (r *DynamoUserRepo) AddBalance(ctx context.Context, userID string, depositAmount float64, updatedAt string) error {
	rec, err := r.GetByUserID(ctx, userID)
	if err != nil || rec == nil {
		if err != nil {
			return err
		}
		return fmt.Errorf("user not found")
	}
	newAmount := rec.Amount + depositAmount
	return r.table.Update("pk", rec.Pk).
		Set("amount", newAmount).
		Set("updated_at", updatedAt).
		Run(ctx)
}
