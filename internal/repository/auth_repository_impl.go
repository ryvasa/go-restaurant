package repository

import (
	"context"
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type AuthRepositoryImpl struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &AuthRepositoryImpl{db}
}

func (r *AuthRepositoryImpl) Login(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	return domain.Auth{}, nil
}
