package service

import (
	"context"
	"errors"
	"github.com/spanwalla/merch-store/internal/entity"
	repomocks "github.com/spanwalla/merch-store/internal/mocks/repository"
	"github.com/spanwalla/merch-store/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestPaymentService_BuyItem(t *testing.T) {
	type args struct {
		ctx   context.Context
		input PaymentBuyItemInput
	}

	type MockBehavior func(u *repomocks.MockUser, i *repomocks.MockItem, o *repomocks.MockOperation, s *repomocks.MockSale, t *repomocks.MockTransactor, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				input: PaymentBuyItemInput{
					UserId:   13,
					ItemName: "hoody",
				},
			},
			mockBehavior: func(u *repomocks.MockUser, i *repomocks.MockItem, o *repomocks.MockOperation, s *repomocks.MockSale, t *repomocks.MockTransactor, args args) {
				fakeItem := entity.Item{
					Id:    10,
					Name:  args.input.ItemName,
					Price: 100,
				}

				i.EXPECT().GetItemByName(args.ctx, args.input.ItemName).Return(fakeItem, nil)

				t.EXPECT().WithinTransaction(args.ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				u.EXPECT().Withdraw(gomock.Any(), args.input.UserId, fakeItem.Price).Return(nil)

				expectedSale := entity.Sale{
					UserId:   args.input.UserId,
					ItemId:   fakeItem.Id,
					Quantity: 1,
				}

				s.EXPECT().Upsert(gomock.Any(), expectedSale).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "item does not exist",
			args: args{
				ctx: context.Background(),
				input: PaymentBuyItemInput{
					UserId:   13,
					ItemName: "bad-item-name",
				},
			},
			mockBehavior: func(u *repomocks.MockUser, i *repomocks.MockItem, o *repomocks.MockOperation, s *repomocks.MockSale, t *repomocks.MockTransactor, args args) {
				i.EXPECT().GetItemByName(args.ctx, args.input.ItemName).Return(entity.Item{}, repository.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "not enough money",
			args: args{
				ctx: context.Background(),
				input: PaymentBuyItemInput{
					UserId:   13,
					ItemName: "powerbank",
				},
			},
			mockBehavior: func(u *repomocks.MockUser, i *repomocks.MockItem, o *repomocks.MockOperation, s *repomocks.MockSale, t *repomocks.MockTransactor, args args) {
				fakeItem := entity.Item{
					Id:    10,
					Name:  args.input.ItemName,
					Price: 200,
				}

				i.EXPECT().GetItemByName(args.ctx, args.input.ItemName).Return(fakeItem, nil)

				t.EXPECT().WithinTransaction(args.ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				u.EXPECT().Withdraw(gomock.Any(), args.input.UserId, fakeItem.Price).Return(repository.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "transaction error",
			args: args{
				ctx: context.Background(),
				input: PaymentBuyItemInput{
					UserId:   13,
					ItemName: "hoody",
				},
			},
			mockBehavior: func(u *repomocks.MockUser, i *repomocks.MockItem, o *repomocks.MockOperation, s *repomocks.MockSale, t *repomocks.MockTransactor, args args) {
				fakeItem := entity.Item{
					Id:    10,
					Name:  args.input.ItemName,
					Price: 100,
				}

				i.EXPECT().GetItemByName(args.ctx, args.input.ItemName).Return(fakeItem, nil)

				t.EXPECT().WithinTransaction(args.ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return errors.New("transaction error")
					})
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := repomocks.NewMockUser(ctrl)
			itemRepo := repomocks.NewMockItem(ctrl)
			operationRepo := repomocks.NewMockOperation(ctrl)
			saleRepo := repomocks.NewMockSale(ctrl)
			transactor := repomocks.NewMockTransactor(ctrl)
			tc.mockBehavior(userRepo, itemRepo, operationRepo, saleRepo, transactor, tc.args)
			s := NewPaymentService(userRepo, itemRepo, operationRepo, saleRepo, transactor)

			err := s.BuyItem(tc.args.ctx, tc.args.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestPaymentService_Transfer(t *testing.T) {
	type args struct {
		ctx   context.Context
		input PaymentTransferInput
	}

	type MockBehavior func(u *repomocks.MockUser, i *repomocks.MockItem, o *repomocks.MockOperation, s *repomocks.MockSale, t *repomocks.MockTransactor, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				input: PaymentTransferInput{
					FromUserId: 13,
					ToUserName: "hoody",
					Amount:     10,
				},
			},
			mockBehavior: func(u *repomocks.MockUser, i *repomocks.MockItem, o *repomocks.MockOperation, s *repomocks.MockSale, t *repomocks.MockTransactor, args args) {
				toUserId := 495
				u.EXPECT().GetUserIdByName(args.ctx, args.input.ToUserName).Return(toUserId, nil)

				t.EXPECT().WithinTransaction(args.ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				u.EXPECT().Withdraw(gomock.Any(), args.input.FromUserId, args.input.Amount).Return(nil)
				u.EXPECT().Deposit(gomock.Any(), toUserId, args.input.Amount).Return(nil)

				expectedOperation := entity.Operation{
					SenderId:   args.input.FromUserId,
					ReceiverId: toUserId,
					Amount:     args.input.Amount,
				}

				o.EXPECT().Upsert(gomock.Any(), expectedOperation).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "not enough coins",
			args: args{
				ctx: context.Background(),
				input: PaymentTransferInput{
					FromUserId: 13,
					ToUserName: "hoody",
					Amount:     1005,
				},
			},
			mockBehavior: func(u *repomocks.MockUser, i *repomocks.MockItem, o *repomocks.MockOperation, s *repomocks.MockSale, t *repomocks.MockTransactor, args args) {
				toUserId := 10039
				u.EXPECT().GetUserIdByName(args.ctx, args.input.ToUserName).Return(toUserId, nil)

				t.EXPECT().WithinTransaction(args.ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				u.EXPECT().Withdraw(gomock.Any(), args.input.FromUserId, args.input.Amount).Return(repository.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "user does not exist",
			args: args{
				ctx: context.Background(),
				input: PaymentTransferInput{
					FromUserId: 13,
					ToUserName: "hoody",
					Amount:     100,
				},
			},
			mockBehavior: func(u *repomocks.MockUser, i *repomocks.MockItem, o *repomocks.MockOperation, s *repomocks.MockSale, t *repomocks.MockTransactor, args args) {
				u.EXPECT().GetUserIdByName(args.ctx, args.input.ToUserName).Return(0, repository.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "self transfer",
			args: args{
				ctx: context.Background(),
				input: PaymentTransferInput{
					FromUserId: 13,
					ToUserName: "hoody",
					Amount:     100,
				},
			},
			mockBehavior: func(u *repomocks.MockUser, i *repomocks.MockItem, o *repomocks.MockOperation, s *repomocks.MockSale, t *repomocks.MockTransactor, args args) {
				toUserId := args.input.FromUserId
				u.EXPECT().GetUserIdByName(args.ctx, args.input.ToUserName).Return(toUserId, nil)
			},
			wantErr: true,
		},
		{
			name: "transaction error",
			args: args{
				ctx: context.Background(),
				input: PaymentTransferInput{
					FromUserId: 13,
					ToUserName: "hoody",
					Amount:     100,
				},
			},
			mockBehavior: func(u *repomocks.MockUser, i *repomocks.MockItem, o *repomocks.MockOperation, s *repomocks.MockSale, t *repomocks.MockTransactor, args args) {
				toUserId := 495
				u.EXPECT().GetUserIdByName(args.ctx, args.input.ToUserName).Return(toUserId, nil)

				t.EXPECT().WithinTransaction(args.ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return errors.New("transaction error")
					})
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := repomocks.NewMockUser(ctrl)
			itemRepo := repomocks.NewMockItem(ctrl)
			operationRepo := repomocks.NewMockOperation(ctrl)
			saleRepo := repomocks.NewMockSale(ctrl)
			transactor := repomocks.NewMockTransactor(ctrl)
			tc.mockBehavior(userRepo, itemRepo, operationRepo, saleRepo, transactor, tc.args)
			s := NewPaymentService(userRepo, itemRepo, operationRepo, saleRepo, transactor)

			err := s.Transfer(tc.args.ctx, tc.args.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
