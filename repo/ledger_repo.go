package repo

import (
	"context"
	"fmt"

	dynamo "github.com/guregu/dynamo/v2"

	"github.com/quangdnptit/booking-be/repomodel"
)

// LedgerRepo writes deposit rows into the same table as user profiles (single-table design).
type LedgerRepo interface {
	AddDeposit(ctx context.Context, rec repomodel.BalanceRecord) error
}

type DynamoLedgerRepo struct {
	table dynamo.Table
}

func NewDynamoLedgerRepo(db *dynamo.DB) *DynamoLedgerRepo {
	return &DynamoLedgerRepo{table: db.Table(TableUsers)}
}

// AddDeposit appends a new row (pk=BALANCE#<user_id>, sk=DEPOSIT#<deposit_id>). All deposits are kept for history.
func (r *DynamoLedgerRepo) AddDeposit(ctx context.Context, rec repomodel.BalanceRecord) error {
	if rec.UserID == "" || rec.DepositID == "" {
		return fmt.Errorf("user_id and deposit_id required")
	}
	rec.Pk = repomodel.PKPrefixBalance + rec.UserID
	rec.Sk = repomodel.SKDepositPrefix + rec.DepositID
	return r.table.Put(rec).Run(ctx)
}
