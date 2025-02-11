package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/internal/repository"
	"github.com/spanwalla/merch-store/pkg/hasher"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int
}

type UserService struct {
	userRepo       repository.User
	passwordHasher hasher.PasswordHasher
	signKey        string
	tokenTTL       time.Duration
}

func NewUserService(userRepo repository.User, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *UserService {
	return &UserService{userRepo: userRepo, passwordHasher: passwordHasher, signKey: signKey, tokenTTL: tokenTTL}
}

func (s *UserService) createUser(ctx context.Context, input UserGenerateTokenInput) (int, error) {
	user := entity.User{
		Name:     input.Name,
		Password: s.passwordHasher.Hash(input.Password),
	}
	userId, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return 0, ErrUserAlreadyExists
		}
		log.Errorf("UserService.createUser - userRepo.CreateUser: %v", err)
		return 0, ErrCannotCreateUser
	}
	return userId, nil
}

func (s *UserService) GenerateToken(ctx context.Context, input UserGenerateTokenInput) (string, error) {
	user, err := s.userRepo.GetUserByNameAndPassword(ctx, input.Name, s.passwordHasher.Hash(input.Password))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			var userId int
			userId, err = s.createUser(ctx, input)
			if err != nil {
				return "", err
			}
			user.Id = userId
		} else {
			log.Errorf("UserService.GenerateToken - userRepo.GetUserByNameAndPassword: %v", err)
			return "", ErrCannotGetUser
		}
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
	})

	// Sign token
	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		log.Errorf("UserService.GenerateToken - token.SignedString: %v", err)
		return "", ErrCannotSignToken
	}

	return tokenString, nil
}

func (s *UserService) VerifyToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return 0, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, ErrCannotParseToken
	}

	return claims.UserId, nil
}
