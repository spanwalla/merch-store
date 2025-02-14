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

func TestOperationRepo_Upsert(t *testing.T) {
	type args struct {
		ctx       context.Context
		operation entity.Operation
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
				operation: entity.Operation{
					SenderId:   1,
					ReceiverId: 49,
					Amount:     100,
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`INSERT INTO operations`).
					WithArgs(args.operation.SenderId, args.operation.ReceiverId, args.operation.Amount).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: false,
		},
		{
			name: "unknown error",
			args: args{
				ctx: context.Background(),
				operation: entity.Operation{
					SenderId:   1,
					ReceiverId: 49,
					Amount:     100,
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`INSERT INTO operations`).
					WithArgs(args.operation.SenderId, args.operation.ReceiverId, args.operation.Amount).
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

			operationRepoMock := NewOperationRepo(postgresMock)

			err := operationRepoMock.Upsert(tc.args.ctx, tc.args.operation)
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
