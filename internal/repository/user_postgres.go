package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) CreateUser(ctx context.Context, user entity.User) (int, error) {
	sql, args, _ := r.Builder.
		Insert("users").
		Columns("name, password").
		Values(user.Name, user.Password).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := r.GetQueryRunner(ctx).QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, ErrAlreadyExists
			}
			return 0, fmt.Errorf("UserRepo.CreateUser - QueryRow: %w", err)
		}
	}

	return id, nil
}

func (r *UserRepo) GetUserByName(ctx context.Context, username string) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("id, name, password, balance").
		From("users").
		Where("name = ?", username).
		ToSql()

	var user entity.User
	err := r.GetQueryRunner(ctx).QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Password,
		&user.Balance,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo.GetUserByName - QueryRow: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserIdByName(ctx context.Context, username string) (int, error) {
	sql, args, _ := r.Builder.
		Select("id").
		From("users").
		Where("name = ?", username).
		ToSql()

	var userId int
	err := r.GetQueryRunner(ctx).QueryRow(ctx, sql, args...).Scan(&userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("UserRepo.GetUserIdByName - QueryRow: %w", err)
	}

	return userId, nil
}

func (r *UserRepo) Withdraw(ctx context.Context, id, amount int) error {
	sql, args, _ := r.Builder.
		Update("users").
		Set("balance", squirrel.Expr("balance - ?", amount)).
		Where(squirrel.And{
			squirrel.Eq{"id": id},
			squirrel.GtOrEq{"balance": amount},
		}).
		ToSql()

	cmdTag, err := r.GetQueryRunner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo.Withdraw - Exec: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *UserRepo) Deposit(ctx context.Context, id, amount int) error {
	sql, args, _ := r.Builder.
		Update("users").
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	cmdTag, err := r.GetQueryRunner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo.Deposit - Exec: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
