package service

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/internal/repository"
)

type PaymentService struct {
	userRepo      repository.User
	itemRepo      repository.Item
	operationRepo repository.Operation
	saleRepo      repository.Sale
	transactor    Transactor
}

func NewPaymentService(userRepo repository.User, itemRepo repository.Item, operationRepo repository.Operation, saleRepo repository.Sale, transactor Transactor) *PaymentService {
	return &PaymentService{
		userRepo:      userRepo,
		itemRepo:      itemRepo,
		operationRepo: operationRepo,
		saleRepo:      saleRepo,
		transactor:    transactor,
	}
}

func (s *PaymentService) Transfer(ctx context.Context, input PaymentTransferInput) error {
	toUserId, err := s.userRepo.GetUserIdByName(ctx, input.ToUserName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		log.Errorf("PaymentService.Transfer - userRepo.GetUserIdByName: %v", err)
		return err
	}

	if toUserId == input.FromUserId {
		return ErrSelfTransfer
	}

	operation := entity.Operation{
		SenderId:   input.FromUserId,
		ReceiverId: toUserId,
		Amount:     input.Amount,
	}

	log.Debug("PaymentService.Transfer: transaction finished")
	return s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		err = s.userRepo.Withdraw(txCtx, operation.SenderId, operation.Amount)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotEnoughBalance
			}
			log.Errorf("PaymentService.Transfer - userRepo.Withdraw: %v", err)
			return ErrCannotTransferCoins
		}

		err = s.userRepo.Deposit(txCtx, operation.ReceiverId, operation.Amount)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrUserNotFound
			}
			log.Errorf("PaymentService.Transfer - userRepo.Deposit: %v", err)
			return ErrCannotTransferCoins
		}

		err = s.operationRepo.Upsert(txCtx, operation)
		if err != nil {
			log.Errorf("PaymentService.Transfer - operationRepo.Upsert: %v", err)
			return ErrCannotTransferCoins
		}

		log.Debug("PaymentService.Transfer: transaction finished")
		return nil
	})
}

func (s *PaymentService) BuyItem(ctx context.Context, input PaymentBuyItemInput) error {
	item, err := s.itemRepo.GetItemByName(ctx, input.ItemName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrItemNotFound
		}
		log.Errorf("PaymentService.BuyItem - itemRepo.GetItemByName: %v", err)
		return ErrCannotBuyItem
	}

	sale := entity.Sale{
		UserId:   input.UserId,
		ItemId:   item.Id,
		Quantity: 1,
	}

	return s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		err = s.userRepo.Withdraw(txCtx, input.UserId, item.Price)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotEnoughBalance
			}
			log.Errorf("PaymentService.BuyItem - userRepo.Withdraw: %v", err)
			return ErrCannotBuyItem
		}

		err = s.saleRepo.Upsert(txCtx, sale)
		if err != nil {
			log.Errorf("PaymentService.BuyItem - saleRepo.Upsert: %v", err)
			return ErrCannotBuyItem
		}

		log.Debug("PaymentService.BuyItem: transaction finished")
		return nil
	})
}
