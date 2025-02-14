package repository

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaleRepo_Upsert(t *testing.T) {
	type args struct {
		ctx  context.Context
		sale entity.Sale
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

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
				sale: entity.Sale{
					UserId:   1,
					ItemId:   10,
					Quantity: 2,
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`INSERT INTO sales`).
					WithArgs(args.sale.UserId, args.sale.ItemId, args.sale.Quantity).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: false,
		},
		{
			name: "unknown error",
			args: args{
				ctx: context.Background(),
				sale: entity.Sale{
					UserId:   1,
					ItemId:   10,
					Quantity: 1,
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`INSERT INTO operations`).
					WithArgs(args.sale.UserId, args.sale.ItemId, args.sale.Quantity).
					WillReturnError(errors.New("some query error"))
			},
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

			saleRepoMock := NewSaleRepo(postgresMock)

			err := saleRepoMock.Upsert(tc.args.ctx, tc.args.sale)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
