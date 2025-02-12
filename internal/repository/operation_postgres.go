package repository

import (
	"context"
	"fmt"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/pkg/postgres"
)

type OperationRepo struct {
	*postgres.Postgres
}

func NewOperationRepo(pg *postgres.Postgres) *OperationRepo {
	return &OperationRepo{pg}
}

func (r *OperationRepo) Upsert(ctx context.Context, operation entity.Operation) error {
	sql, args, _ := r.Builder.
		Insert("operations").
		Columns("sender_id, receiver_id, amount").
		Values(operation.SenderId, operation.ReceiverId, operation.Amount).
		Suffix("ON CONFLICT (sender_id, receiver_id) DO UPDATE SET amount = operations.amount + EXCLUDED.amount").
		ToSql()

	_, err := r.GetQueryRunner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OperationRepo.Upsert - Exec: %w", err)
	}

	return nil
}
