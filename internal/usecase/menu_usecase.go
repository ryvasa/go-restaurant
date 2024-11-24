package usecase

import (
	"context"
	"mime/multipart"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type MenuUsecase interface {
	GetAll(ctx context.Context) ([]domain.Menu, error)
	Create(ctx context.Context, req dto.CreateMenuRequest, file multipart.File) (domain.Menu, error)
	Get(ctx context.Context, id string) (domain.Menu, error)
	Update(ctx context.Context, id string, req dto.UpdateMenuRequest, file multipart.File) (domain.Menu, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) (domain.Menu, error)
}
