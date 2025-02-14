package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserReportRepo_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

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
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				expectedInventory := []entity.Inventory{
					{
						Type:     "cup",
						Quantity: 0,
					},
					{
						Type:     "book",
						Quantity: 4,
					},
				}
				expectedCoinHistory := entity.CoinHistory{
					Received: []entity.ReceivedTransaction{
						{
							FromUser: "user1",
							Amount:   10,
						},
						{
							FromUser: "user2",
							Amount:   20,
						},
					},
					Sent: []entity.SentTransaction{
						{
							ToUser: "user4",
							Amount: 42,
						},
						{
							ToUser: "user2",
							Amount: 5,
						},
					},
				}

				expectedInventoryJSON, _ := json.Marshal(expectedInventory)
				expectedCoinHistoryJSON, _ := json.Marshal(expectedCoinHistory)
				rows := pgxmock.NewRows([]string{"balance", "inventory", "coin_history"}).
					AddRow(100, expectedInventoryJSON, expectedCoinHistoryJSON)

				m.ExpectQuery(`SELECT u.balance`).
					WithArgs(args.id).
					WillReturnRows(rows)
			},
			want: entity.UserReport{
				Coins: 100,
				Inventory: []entity.Inventory{
					{
						Type:     "cup",
						Quantity: 0,
					},
					{
						Type:     "book",
						Quantity: 4,
					},
				},
				CoinHistory: entity.CoinHistory{
					Received: []entity.ReceivedTransaction{
						{
							FromUser: "user1",
							Amount:   10,
						},
						{
							FromUser: "user2",
							Amount:   20,
						},
					},
					Sent: []entity.SentTransaction{
						{
							ToUser: "user4",
							Amount: 42,
						},
						{
							ToUser: "user2",
							Amount: 5,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "user not found",
			args: args{
				ctx: context.Background(),
				id:  130,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery(`SELECT u.balance`).
					WithArgs(args.id).
					WillReturnError(pgx.ErrNoRows)
			},
			want:    entity.UserReport{},
			wantErr: true,
		},
		{
			name: "unknown error",
			args: args{
				ctx: context.Background(),
				id:  130,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery(`SELECT u.balance`).
					WithArgs(args.id).
					WillReturnError(errors.New("unexpected error"))
			},
			want:    entity.UserReport{},
			wantErr: true,
		},
		{
			name: "corrupted inventory json from db",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				expectedInventory := []entity.Inventory{
					{
						Type:     "cup",
						Quantity: 0,
					},
					{
						Type:     "book",
						Quantity: 4,
					},
				}
				expectedCoinHistory := entity.CoinHistory{
					Received: []entity.ReceivedTransaction{
						{
							FromUser: "user1",
							Amount:   10,
						},
						{
							FromUser: "user2",
							Amount:   20,
						},
					},
					Sent: []entity.SentTransaction{
						{
							ToUser: "user4",
							Amount: 42,
						},
						{
							ToUser: "user2",
							Amount: 5,
						},
					},
				}

				expectedInventoryJSON, _ := json.Marshal(expectedInventory)
				expectedCoinHistoryJSON, _ := json.Marshal(expectedCoinHistory)
				rows := pgxmock.NewRows([]string{"balance", "inventory", "coin_history"}).
					AddRow(100, append(expectedInventoryJSON, '1'), expectedCoinHistoryJSON)

				m.ExpectQuery(`SELECT u.balance`).
					WithArgs(args.id).
					WillReturnRows(rows)
			},
			want:    entity.UserReport{},
			wantErr: true,
		},
		{
			name: "corrupted history json from db",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				expectedInventory := []entity.Inventory{
					{
						Type:     "cup",
						Quantity: 0,
					},
					{
						Type:     "book",
						Quantity: 4,
					},
				}
				expectedCoinHistory := entity.CoinHistory{
					Received: []entity.ReceivedTransaction{
						{
							FromUser: "user1",
							Amount:   10,
						},
						{
							FromUser: "user2",
							Amount:   20,
						},
					},
					Sent: []entity.SentTransaction{
						{
							ToUser: "user4",
							Amount: 42,
						},
						{
							ToUser: "user2",
							Amount: 5,
						},
					},
				}

				expectedInventoryJSON, _ := json.Marshal(expectedInventory)
				expectedCoinHistoryJSON, _ := json.Marshal(expectedCoinHistory)
				rows := pgxmock.NewRows([]string{"balance", "inventory", "coin_history"}).
					AddRow(100, expectedInventoryJSON, append(expectedCoinHistoryJSON, '!'))

				m.ExpectQuery(`SELECT u.balance`).
					WithArgs(args.id).
					WillReturnRows(rows)
			},
			want:    entity.UserReport{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}

			userReportRepoMock := NewUserReportRepo(postgresMock)

			got, err := userReportRepoMock.Get(tc.args.ctx, tc.args.id)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
