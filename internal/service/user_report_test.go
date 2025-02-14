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
						Inventory: []struct {
							Type     string `json:"type"`
							Quantity int    `json:"quantity"`
						}{
							{
								Type:     "t-shirt",
								Quantity: 10,
							},
							{
								Type:     "powerbank",
								Quantity: 0,
							},
						},
						CoinHistory: struct {
							Received []struct {
								FromUser string `json:"fromUser"`
								Amount   int    `json:"amount"`
							} `json:"received"`
							Sent []struct {
								ToUser string `json:"toUser"`
								Amount int    `json:"amount"`
							} `json:"sent"`
						}{
							Received: []struct {
								FromUser string `json:"fromUser"`
								Amount   int    `json:"amount"`
							}{
								{
									FromUser: "test-user-1",
									Amount:   10,
								},
								{
									FromUser: "test-user-2",
									Amount:   20,
								},
							},
							Sent: []struct {
								ToUser string `json:"toUser"`
								Amount int    `json:"amount"`
							}{},
						},
					}, nil)
			},
			want: entity.UserReport{
				Coins: 1000,
				Inventory: []struct {
					Type     string `json:"type"`
					Quantity int    `json:"quantity"`
				}{
					{
						Type:     "t-shirt",
						Quantity: 10,
					},
					{
						Type:     "powerbank",
						Quantity: 0,
					},
				},
				CoinHistory: struct {
					Received []struct {
						FromUser string `json:"fromUser"`
						Amount   int    `json:"amount"`
					} `json:"received"`
					Sent []struct {
						ToUser string `json:"toUser"`
						Amount int    `json:"amount"`
					} `json:"sent"`
				}{
					Received: []struct {
						FromUser string `json:"fromUser"`
						Amount   int    `json:"amount"`
					}{
						{
							FromUser: "test-user-1",
							Amount:   10,
						},
						{
							FromUser: "test-user-2",
							Amount:   20,
						},
					},
					Sent: []struct {
						ToUser string `json:"toUser"`
						Amount int    `json:"amount"`
					}{},
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
