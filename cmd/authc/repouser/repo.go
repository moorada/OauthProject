package repouser

import (
	"context"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (user UserInfo, err error)
	GetUserById(ctx context.Context, id string) (user UserInfo, err error)
}

type UserInfo struct {
	ID       string
	Email    string
	Nome     string
	Cognome  string
	Password string
}
