package repository

import (
	"context"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/pkg/postgres"
)

//go:generate mockgen -source=repository.go -destination=../mocks/repository/mock.go -package=repomocks

// Transactor определяет интерфейс для работы с транзакциями.
type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Operation interface {
	Upsert(ctx context.Context, operation entity.Operation) error
}

type Item interface {
	GetItemByName(ctx context.Context, name string) (entity.Item, error)
}

type Sale interface {
	Upsert(ctx context.Context, sale entity.Sale) error
}

type User interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserByName(ctx context.Context, username string) (entity.User, error)
	GetUserIdByName(ctx context.Context, username string) (int, error)
	Withdraw(ctx context.Context, id, amount int) error
	Deposit(ctx context.Context, id, amount int) error
}

type UserReport interface {
	Get(ctx context.Context, id int) (entity.UserReport, error)
}

type Repositories struct {
	Operation
	Item
	Sale
	User
	UserReport
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Operation:  NewOperationRepo(pg),
		Item:       NewItemRepo(pg),
		Sale:       NewSaleRepo(pg),
		User:       NewUserRepo(pg),
		UserReport: NewUserReportRepo(pg),
	}
}
