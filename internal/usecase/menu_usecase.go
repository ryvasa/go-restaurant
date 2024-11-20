package usecase

import (
	"context"
	"time"

	"github.com/ryvasa/go-restaurant/internal/domain"
)

type menuUsecase struct {
	menuRepo domain.MenuRepository
}

func NewMenuUsecase(menuRepo domain.MenuRepository) domain.MenuUsecase {
	return &menuUsecase{menuRepo}
}
func (u *menuUsecase) GetAll(ctx context.Context) ([]domain.Menu, error) {
	return u.menuRepo.GetAll(ctx)
}

func (u *menuUsecase) Create(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return u.menuRepo.Create(ctx, menu)
}

func (u *menuUsecase) Get(ctx context.Context, id string) (domain.Menu, error) {
	return u.menuRepo.Get(ctx, id)
}

func (u *menuUsecase) Update(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	return u.menuRepo.Update(ctx, menu)
}

func (u *menuUsecase) Delete(ctx context.Context, id string) error {
	return u.menuRepo.Delete(ctx, id)
}
