package repository

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type AuthRepository interface {
	Login(ctx context.Context, auth domain.Auth) (domain.Auth, error)
}
