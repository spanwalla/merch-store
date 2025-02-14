package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/pkg/postgres"
)

type UserReportRepo struct {
	*postgres.Postgres
}

func NewUserReportRepo(pg *postgres.Postgres) *UserReportRepo {
	return &UserReportRepo{pg}
}

func (r *UserReportRepo) Get(ctx context.Context, id int) (entity.UserReport, error) {
	inventorySubquery := r.Builder.
		Select("COALESCE(jsonb_agg(jsonb_build_object('type', i.name, 'quantity', COALESCE(s.quantity, 0))), '[]'::jsonb)").
		From("items i").
		LeftJoin("sales s ON i.id = s.item_id AND s.user_id = u.id")

	sentSubquery := r.Builder.
		Select("jsonb_agg(jsonb_build_object('toUser', r.name, 'amount', o.amount))").
		From("operations o").
		Join("users r ON o.receiver_id = r.id").
		Where("o.sender_id = u.id")

	receivedSubquery := r.Builder.
		Select("jsonb_agg(jsonb_build_object('fromUser', s.name, 'amount', o.amount))").
		From("operations o").
		Join("users s ON o.sender_id = s.id").
		Where("o.receiver_id = u.id")

	inventorySql, _, _ := squirrel.Expr("(?) AS inventory", inventorySubquery).ToSql()
	historySql, _, _ := squirrel.Expr("jsonb_build_object('sent', COALESCE((?), '[]'::jsonb), 'received', COALESCE((?), '[]'::jsonb)) AS coin_history", sentSubquery, receivedSubquery).ToSql()

	sql, args, _ := r.Builder.
		Select(
			"u.balance",
			inventorySql,
			historySql,
		).
		From("users u").
		Where("u.id = ?", id).ToSql()

	var userReport entity.UserReport
	var inventoryJSON []byte
	var historyJSON []byte
	err := r.GetQueryRunner(ctx).QueryRow(ctx, sql, args...).Scan(
		&userReport.Coins,
		&inventoryJSON,
		&historyJSON,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserReport{}, ErrNotFound
		}
		return entity.UserReport{}, fmt.Errorf("UserReportRepo.Get - QueryRow: %w", err)
	}

	err = json.Unmarshal(inventoryJSON, &userReport.Inventory)
	if err != nil {
		return entity.UserReport{}, fmt.Errorf("UserReportRepo.Get - Unmarshal Inventory: %w", err)
	}
	err = json.Unmarshal(historyJSON, &userReport.CoinHistory)
	if err != nil {
		return entity.UserReport{}, fmt.Errorf("UserReportRepo.Get - Unmarshal History: %w", err)
	}

	return userReport, nil
}
