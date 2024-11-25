package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type AuthUsecase interface {
	Login(ctx context.Context, req dto.LoginDto) (domain.Auth, error)
}
