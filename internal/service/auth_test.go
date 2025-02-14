package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/spanwalla/merch-store/internal/entity"
	hashermocks "github.com/spanwalla/merch-store/internal/mocks/hasher"
	"github.com/spanwalla/merch-store/internal/mocks/repository"
	"github.com/spanwalla/merch-store/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestAuthService_createUser(t *testing.T) {
	const secret = "jwt_test_secret"
	const tokenTTL = 2 * time.Hour

	type args struct {
		ctx   context.Context
		input AuthGenerateTokenInput
	}

	type MockBehavior func(u *repomocks.MockUser, h *hashermocks.MockPasswordHasher, s string, ttl time.Duration, args args)

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
				input: AuthGenerateTokenInput{
					Name:     "marcus-web-designer",
					Password: "simplePa66!",
				},
			},
			mockBehavior: func(u *repomocks.MockUser, h *hashermocks.MockPasswordHasher, s string, ttl time.Duration, args args) {
				h.EXPECT().Hash(args.input.Password).
					Return(args.input.Password)
				u.EXPECT().CreateUser(args.ctx, entity.User{Name: args.input.Name, Password: args.input.Password}).
					Return(1, nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "already exists",
			args: args{
				ctx: context.Background(),
				input: AuthGenerateTokenInput{
					Name:     "marcus-web-designer",
					Password: "simplePa66!",
				},
			},
			mockBehavior: func(u *repomocks.MockUser, h *hashermocks.MockPasswordHasher, s string, ttl time.Duration, args args) {
				h.EXPECT().Hash(args.input.Password).
					Return(args.input.Password)
				u.EXPECT().CreateUser(args.ctx, entity.User{Name: args.input.Name, Password: args.input.Password}).
					Return(0, repository.ErrAlreadyExists)
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := repomocks.NewMockUser(ctrl)
			hasher := hashermocks.NewMockPasswordHasher(ctrl)
			tc.mockBehavior(userRepo, hasher, secret, tokenTTL, tc.args)

			s := NewAuthService(userRepo, hasher, secret, tokenTTL)

			got, err := s.createUser(tc.args.ctx, tc.args.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	const secret = "jwt_test_secret"
	const tokenTTL = 2 * time.Hour

	type args struct {
		ctx   context.Context
		input AuthGenerateTokenInput
	}

	type MockBehavior func(u *repomocks.MockUser, h *hashermocks.MockPasswordHasher, s string, ttl time.Duration, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      bool
	}{
		{
			name: "registration success",
			args: args{
				ctx: context.Background(),
				input: AuthGenerateTokenInput{
					Name:     "marcus-web-designer",
					Password: "simplePa66!",
				},
			},
			mockBehavior: func(u *repomocks.MockUser, h *hashermocks.MockPasswordHasher, s string, ttl time.Duration, args args) {
				u.EXPECT().GetUserByName(args.ctx, args.input.Name).
					Return(entity.User{}, repository.ErrNotFound)
				h.EXPECT().Hash(args.input.Password).
					Return(args.input.Password)
				u.EXPECT().CreateUser(args.ctx, entity.User{Name: args.input.Name, Password: args.input.Password}).
					Return(1, nil)
			},
			wantErr: false,
		},
		{
			name: "authorization success",
			args: args{
				ctx: context.Background(),
				input: AuthGenerateTokenInput{
					Name:     "marcus-web-designer",
					Password: "simplePa66!",
				},
			},
			mockBehavior: func(u *repomocks.MockUser, h *hashermocks.MockPasswordHasher, s string, ttl time.Duration, args args) {
				u.EXPECT().GetUserByName(args.ctx, args.input.Name).
					Return(entity.User{Id: 1, Name: args.input.Name, Password: args.input.Password}, nil)
				h.EXPECT().Hash(args.input.Password).
					Return(args.input.Password)
			},
			wantErr: false,
		},
		{
			name: "authorization failed due to wrong password",
			args: args{
				ctx: context.Background(),
				input: AuthGenerateTokenInput{
					Name:     "marcus-web-designer",
					Password: "simplePa66!",
				},
			},
			mockBehavior: func(u *repomocks.MockUser, h *hashermocks.MockPasswordHasher, s string, ttl time.Duration, args args) {
				u.EXPECT().GetUserByName(args.ctx, args.input.Name).
					Return(entity.User{Id: 1, Name: args.input.Name, Password: "another-password"}, nil)
				h.EXPECT().Hash(args.input.Password).
					Return(args.input.Password)
			},
			wantErr: true,
		},
		{
			name: "get user failed for unknown reason",
			args: args{
				ctx: context.Background(),
				input: AuthGenerateTokenInput{
					Name:     "marcus-web-designer",
					Password: "simplePa66!",
				},
			},
			mockBehavior: func(u *repomocks.MockUser, h *hashermocks.MockPasswordHasher, s string, ttl time.Duration, args args) {
				u.EXPECT().GetUserByName(args.ctx, args.input.Name).
					Return(entity.User{}, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := repomocks.NewMockUser(ctrl)
			hasher := hashermocks.NewMockPasswordHasher(ctrl)
			tc.mockBehavior(userRepo, hasher, secret, tokenTTL, tc.args)

			s := NewAuthService(userRepo, hasher, secret, tokenTTL)

			got, err := s.GenerateToken(tc.args.ctx, tc.args.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}

func TestAuthService_VerifyToken(t *testing.T) {
	const secret = "jwt_test_secret"
	const wrongSecret = "wrong_secret"
	const userId = 17
	const tokenTTL = 2 * time.Hour

	type args struct {
		tokenString string
	}

	generateJwt := func(id int, issuedAt time.Time, signKey string) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: issuedAt.Add(tokenTTL).Unix(),
				IssuedAt:  issuedAt.Unix(),
			},
			UserId: userId,
		})

		tokenString, _ := token.SignedString([]byte(signKey))
		return tokenString
	}

	testCases := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				tokenString: generateJwt(userId, time.Now(), secret),
			},
			want:    userId,
			wantErr: false,
		},
		{
			name: "expired token",
			args: args{
				tokenString: generateJwt(userId, time.Now().Add(-10*time.Hour), secret),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "wrong secret",
			args: args{
				tokenString: generateJwt(userId, time.Now(), wrongSecret),
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := repomocks.NewMockUser(ctrl)
			hasher := hashermocks.NewMockPasswordHasher(ctrl)
			s := NewAuthService(userRepo, hasher, secret, tokenTTL)

			got, err := s.VerifyToken(tc.args.tokenString)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
