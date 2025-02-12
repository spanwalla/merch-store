package service

import "errors"

var (
	ErrCannotSignToken  = errors.New("cannot sign token")
	ErrCannotParseToken = errors.New("cannot parse token")

	ErrWrongPassword     = errors.New("wrong password")
	ErrCannotGetUser     = errors.New("cannot get user")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCannotCreateUser  = errors.New("cannot create user")

	ErrNotEnoughBalance    = errors.New("not enough balance")
	ErrItemNotFound        = errors.New("item not found")
	ErrCannotBuyItem       = errors.New("cannot buy item")
	ErrUserNotFound        = errors.New("user not found")
	ErrCannotTransferCoins = errors.New("cannot transfer coins")
	ErrSelfTransfer        = errors.New("cannot transfer coins to yourself")
)
