package service

import (
	"context"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/internal/repository"
	"github.com/spanwalla/merch-store/pkg/hasher"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go -package=service_mocks

// Transactor определяет интерфейс для работы с транзакциями.
type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type AuthGenerateTokenInput struct {
	Name     string
	Password string
}

type Auth interface {
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	VerifyToken(tokenString string) (int, error)
}

type PaymentTransferInput struct {
	FromUserId int
	ToUserName string
	Amount     int
}

type PaymentBuyItemInput struct {
	UserId   int
	ItemName string
}

type Payment interface {
	Transfer(ctx context.Context, input PaymentTransferInput) error
	BuyItem(ctx context.Context, input PaymentBuyItemInput) error
}

type UserReport interface {
	Get(ctx context.Context, userId int) (entity.UserReport, error)
}

type Services struct {
	Auth
	Payment
	UserReport
}

type Dependencies struct {
	Repos      *repository.Repositories
	Hasher     hasher.PasswordHasher
	SignKey    string
	TokenTTL   time.Duration
	Transactor Transactor
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		Auth:       NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
		Payment:    NewPaymentService(deps.Repos.User, deps.Repos.Item, deps.Repos.Operation, deps.Repos.Sale, deps.Transactor),
		UserReport: NewUserReportService(deps.Repos.UserReport),
	}
}
