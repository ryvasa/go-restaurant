package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/domain"
	"github.com/ryvasa/go-restaurant/internal/repository"
)

// TODO: update move bisnis logic to user usecase
type UserUsecaseImpl struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{userRepo}
}

func (u *UserUsecaseImpl) GetAll(ctx context.Context) ([]domain.User, error) {
	return u.userRepo.GetAll(ctx)
}

func (u *UserUsecaseImpl) Create(ctx context.Context, user domain.User) (domain.User, error) {
	return u.userRepo.Create(ctx, user)
}

func (u *UserUsecaseImpl) Get(ctx context.Context, id string) (domain.User, error) {
	return u.userRepo.Get(ctx, id)
}

func (u *UserUsecaseImpl) Update(ctx context.Context, user domain.User) (domain.User, error) {
	return u.userRepo.Update(ctx, user)
}
