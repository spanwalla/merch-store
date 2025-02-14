package repository

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepo_CreateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user entity.User
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         int
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Name:     "testUserMorty",
					Password: "hisSuperPassword",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(1)

				m.ExpectQuery(`INSERT INTO users`).
					WithArgs(args.user.Name, args.user.Password).
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "user already exists",
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Name:     "existentUser",
					Password: "hisSuperPassword",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`INSERT INTO users`).
					WithArgs(args.user.Name, args.user.Password).
					WillReturnError(&pgconn.PgError{Code: "23505"})
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "unknown error",
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Name:     "normalUser",
					Password: "hisSuperPassword",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`INSERT INTO users`).
					WithArgs(args.user.Name, args.user.Password).
					WillReturnError(errors.New("unexpected error"))
			},
			want:    0,
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

			userRepoMock := NewUserRepo(postgresMock)

			got, err := userRepoMock.CreateUser(tc.args.ctx, tc.args.user)
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

func TestUserRepo_GetUserByName(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         entity.User
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				username: "test-name",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id", "name", "password", "balance"}).
					AddRow(1, "test-name", "password", 1000)

				m.ExpectQuery(`SELECT id, name, password, balance`).
					WithArgs(args.username).
					WillReturnRows(rows)
			},
			want: entity.User{
				Id:       1,
				Name:     "test-name",
				Password: "password",
				Balance:  1000,
			},
			wantErr: false,
		},
		{
			name: "user not found",
			args: args{
				ctx:      context.Background(),
				username: "wrongUsername",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery(`SELECT id, name, password, balance`).
					WithArgs(args.username).
					WillReturnError(pgx.ErrNoRows)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "unknown error",
			args: args{
				ctx:      context.Background(),
				username: "normalName1",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery(`SELECT id, name, password, balance`).
					WithArgs(args.username).
					WillReturnError(errors.New("unexpected error"))
			},
			want:    entity.User{},
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

			userRepoMock := NewUserRepo(postgresMock)

			got, err := userRepoMock.GetUserByName(tc.args.ctx, tc.args.username)
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

func TestUserRepo_GetUserIdByName(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         int
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				username: "test",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(10)

				m.ExpectQuery(`SELECT id`).
					WithArgs(args.username).
					WillReturnRows(rows)
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "user not found",
			args: args{
				ctx:      context.Background(),
				username: "notFoundUser",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery(`SELECT id`).
					WithArgs(args.username).
					WillReturnError(pgx.ErrNoRows)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "unknown error",
			args: args{
				ctx:      context.Background(),
				username: "prettyThings",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery(`SELECT id`).
					WithArgs(args.username).
					WillReturnError(errors.New("unexpected error"))
			},
			want:    0,
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

			userRepoMock := NewUserRepo(postgresMock)

			got, err := userRepoMock.GetUserIdByName(tc.args.ctx, tc.args.username)
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

func TestUserRepo_Withdraw(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     int
		amount int
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
				ctx:    context.Background(),
				id:     1,
				amount: 100,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`UPDATE users`).
					WithArgs(args.amount, args.id, args.amount).
					WillReturnResult(pgxmock.NewResult(`UPDATE`, 1))
			},
			wantErr: false,
		},
		{
			name: "no rows affected",
			args: args{
				ctx:    context.Background(),
				id:     1,
				amount: 100,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`UPDATE users`).
					WithArgs(args.amount, args.id, args.amount).
					WillReturnResult(pgxmock.NewResult(`UPDATE`, 0))
			},
			wantErr: true,
		},
		{
			name: "unknown error",
			args: args{
				ctx:    context.Background(),
				id:     1,
				amount: 100,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`UPDATE users`).
					WithArgs(args.amount, args.id, args.amount).
					WillReturnError(errors.New("unexpected error"))
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

			userRepoMock := NewUserRepo(postgresMock)

			err := userRepoMock.Withdraw(tc.args.ctx, tc.args.id, tc.args.amount)
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

func TestUserRepo_Deposit(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     int
		amount int
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
				ctx:    context.Background(),
				id:     1,
				amount: 100,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`UPDATE users`).
					WithArgs(args.amount, args.id).
					WillReturnResult(pgxmock.NewResult(`UPDATE`, 1))
			},
			wantErr: false,
		},
		{
			name: "unknown error",
			args: args{
				ctx:    context.Background(),
				id:     1,
				amount: 100,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`UPDATE users`).
					WithArgs(args.amount, args.id).
					WillReturnError(errors.New("unexpected error"))
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

			userRepoMock := NewUserRepo(postgresMock)

			err := userRepoMock.Deposit(tc.args.ctx, tc.args.id, tc.args.amount)
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
