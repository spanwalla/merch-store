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

func TestUserReportService_Get(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId int
	}

	type MockBehavior func(u *repomocks.MockUserReport, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         entity.UserReport
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userId: 49,
			},
			mockBehavior: func(u *repomocks.MockUserReport, args args) {
				u.EXPECT().Get(args.ctx, args.userId).
					Return(entity.UserReport{
						Coins: 1000,
						Inventory: []entity.Inventory{
							{
								Type:     "t-shirt",
								Quantity: 10,
							},
							{
								Type:     "powerbank",
								Quantity: 0,
							},
						},
						CoinHistory: entity.CoinHistory{
							Received: []entity.ReceivedTransaction{
								{
									FromUser: "test-user-1",
									Amount:   10,
								},
								{
									FromUser: "test-user-2",
									Amount:   20,
								},
							},
							Sent: []entity.SentTransaction{},
						},
					}, nil)
			},
			want: entity.UserReport{
				Coins: 1000,
				Inventory: []entity.Inventory{
					{
						Type:     "t-shirt",
						Quantity: 10,
					},
					{
						Type:     "powerbank",
						Quantity: 0,
					},
				},
				CoinHistory: entity.CoinHistory{
					Received: []entity.ReceivedTransaction{
						{
							FromUser: "test-user-1",
							Amount:   10,
						},
						{
							FromUser: "test-user-2",
							Amount:   20,
						},
					},
					Sent: []entity.SentTransaction{},
				},
			},
			wantErr: false,
		},
		{
			name: "user not found",
			args: args{
				ctx:    context.Background(),
				userId: 20,
			},
			mockBehavior: func(u *repomocks.MockUserReport, args args) {
				u.EXPECT().Get(args.ctx, args.userId).
					Return(entity.UserReport{}, repository.ErrNotFound)
			},
			want:    entity.UserReport{},
			wantErr: true,
		},
		{
			name: "some error from repository",
			args: args{
				ctx:    context.Background(),
				userId: -20,
			},
			mockBehavior: func(u *repomocks.MockUserReport, args args) {
				u.EXPECT().Get(args.ctx, args.userId).
					Return(entity.UserReport{}, errors.New("some error"))
			},
			want:    entity.UserReport{},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userReportRepo := repomocks.NewMockUserReport(ctrl)
			tc.mockBehavior(userReportRepo, tc.args)

			s := NewUserReportService(userReportRepo)

			got, err := s.Get(tc.args.ctx, tc.args.userId)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
