package domain

import "errors"

var (
	ErrEntityNotFound      = errors.New("error entity not found")
	ErrEntityAlreadyExists = errors.New("error entity already exists")
	ErrInvalidCredentails  = errors.New("error invalid credentails")
	ErrInvalidToken        = errors.New("error invalid token")
)
