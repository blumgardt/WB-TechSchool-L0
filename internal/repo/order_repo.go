package repo

import (
	"WB-TechSchool-L0/internal/domain"
	"context"
)

type OrderRepo interface {
	GetById(ctx context.Context, uid string) (*domain.Order, error)
	GetAll(ctx context.Context) ([]domain.Order, error)
	Save(ctx context.Context, order *domain.Order) error
}
