package repository

import (
	"context"
	"fmt"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/pkg/postgres"
)

type SaleRepo struct {
	*postgres.Postgres
}

func NewSaleRepo(pg *postgres.Postgres) *SaleRepo {
	return &SaleRepo{pg}
}

func (r *SaleRepo) Upsert(ctx context.Context, sale entity.Sale) error {
	sql, args, _ := r.Builder.
		Insert("sales").
		Columns("user_id, item_id, quantity").
		Values(sale.UserId, sale.ItemId, sale.Quantity).
		Suffix("ON CONFLICT (user_id, item_id) DO UPDATE SET quantity = sales.quantity + EXCLUDED.quantity").
		ToSql()

	_, err := r.GetQueryRunner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SaleRepo.Upsert - Exec: %w", err)
	}

	return nil
}
