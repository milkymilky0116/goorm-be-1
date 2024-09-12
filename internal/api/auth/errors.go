package auth

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("password is incorrect")
	ErrGenerateTokenFail = errors.New("generate token failed")
	ErrTransactionFail   = errors.New("db transaction failed")
)
