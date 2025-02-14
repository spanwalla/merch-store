package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/pkg/postgres"
)

type ItemRepo struct {
	*postgres.Postgres
}

func NewItemRepo(pg *postgres.Postgres) *ItemRepo {
	return &ItemRepo{pg}
}

func (r *ItemRepo) GetItemByName(ctx context.Context, name string) (entity.Item, error) {
	sql, args, _ := r.Builder.
		Select("id, name, price").
		From("items").
		Where("name = ?", name).
		ToSql()

	var item entity.Item
	err := r.GetQueryRunner(ctx).QueryRow(ctx, sql, args...).Scan(
		&item.Id,
		&item.Name,
		&item.Price,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Item{}, ErrNotFound
		}
		return entity.Item{}, fmt.Errorf("ItemRepo.GetItemByName - QueryRow: %w", err)
	}

	return item, nil
}
