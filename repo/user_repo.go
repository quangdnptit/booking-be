package repo

import (
	"context"
	"fmt"
	"strings"

	dynamo "github.com/guregu/dynamo/v2"

	"booking-be/repomodel"
)

const TableUsers = "users"

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*repomodel.UserRecord, error)
	Create(ctx context.Context, rec repomodel.UserRecord) error
	UpdateAudit(ctx context.Context, email, updatedAt string) error
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

func (r *DynamoUserRepo) GetByEmail(ctx context.Context, email string) (*repomodel.UserRecord, error) {
	key := normalizeEmail(email)
	if key == "" {
		return nil, nil
	}
	var rec repomodel.UserRecord
	err := r.table.Get("email", key).One(ctx, &rec)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &rec, nil
}

func (r *DynamoUserRepo) Create(ctx context.Context, rec repomodel.UserRecord) error {
	rec.Email = normalizeEmail(rec.Email)
	return r.table.Put(rec).If("attribute_not_exists(email)").Run(ctx)
}

func (r *DynamoUserRepo) UpdateAudit(ctx context.Context, email, updatedAt string) error {
	return r.table.Update("email", normalizeEmail(email)).Set("updated_at", updatedAt).Run(ctx)
}
