package usecase

import (
	"github.com/ryvasa/go-restaurant/internal/domain"
)

type menuUsecase struct {
	menuRepo domain.MenuRepository
}

func NewMenuUsecase(menuRepo domain.MenuRepository) domain.MenuUsecase {
	return &menuUsecase{menuRepo}
}

func (u *menuUsecase) GetAll() ([]domain.Menu, error) {
	return u.menuRepo.GetAll()
}

func (u *menuUsecase) Create(menu domain.Menu) (domain.Menu, error) {
	return u.menuRepo.Create(menu)
}

// func (u *menuUsecase) GetByID(id int) (domain.Menu, error) {
// 	return u.menuRepo.GetByID(id)
// }

// func (u *menuUsecase) Update(menu domain.Menu) error {
// 	return u.menuRepo.Update(menu)
// }

// func (u *menuUsecase) Delete(id int) error {
// 	return u.menuRepo.Delete(id)
// }
