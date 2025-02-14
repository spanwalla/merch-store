package repository

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItemRepo_GetItemByName(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         entity.Item
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				name: "sweater",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id", "name", "price"}).
					AddRow(1, args.name, 100)

				m.ExpectQuery(`SELECT id, name, price`).
					WithArgs(args.name).
					WillReturnRows(rows)
			},
			want: entity.Item{
				Id:    1,
				Name:  "sweater",
				Price: 100,
			},
			wantErr: false,
		},
		{
			name: "unknown item",
			args: args{
				ctx:  context.Background(),
				name: "unknown",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery(`SELECT id, name, price`).
					WithArgs(args.name).
					WillReturnError(pgx.ErrNoRows)
			},
			want:    entity.Item{},
			wantErr: true,
		},
		{
			name: "unknown error",
			args: args{
				ctx:  context.Background(),
				name: "something",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery(`SELECT id, name, price`).
					WithArgs(args.name).
					WillReturnError(errors.New("some query error"))
			},
			want:    entity.Item{},
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

			itemRepoMock := NewItemRepo(postgresMock)

			got, err := itemRepoMock.GetItemByName(tc.args.ctx, tc.args.name)
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
