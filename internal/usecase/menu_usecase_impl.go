package usecase

import (
	"context"
	"time"

	"github.com/ryvasa/go-restaurant/internal/domain"
	"github.com/ryvasa/go-restaurant/internal/repository"
)

type MenuUsecaseImpl struct {
	menuRepo repository.MenuRepository
}

func NewMenuUsecase(menuRepo repository.MenuRepository) MenuUsecase {
	return &MenuUsecaseImpl{menuRepo}
}
func (u *MenuUsecaseImpl) GetAll(ctx context.Context) ([]domain.Menu, error) {
	return u.menuRepo.GetAll(ctx)
}

func (u *MenuUsecaseImpl) Create(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return u.menuRepo.Create(ctx, menu)
}

func (u *MenuUsecaseImpl) Get(ctx context.Context, id string) (domain.Menu, error) {
	return u.menuRepo.Get(ctx, id)
}

func (u *MenuUsecaseImpl) Update(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	return u.menuRepo.Update(ctx, menu)
}

func (u *MenuUsecaseImpl) Delete(ctx context.Context, id string) error {
	return u.menuRepo.Delete(ctx, id)
}
