package service

import (
	"WB-TechSchool-L0/internal/domain"
	"WB-TechSchool-L0/internal/repo"
	"WB-TechSchool-L0/internal/service/redis"
	"context"
	"log"
)

type OrderService struct {
	repo        repo.OrderRepo
	redisClient *redis.Client
}

func NewOrderService(repo repo.OrderRepo, client *redis.Client) *OrderService {
	return &OrderService{
		repo:        repo,
		redisClient: client,
	}
}

func (r *OrderService) GetOrderById(id string, ctx context.Context) (*domain.Order, error) {
	order, err := r.redisClient.GetOrder(ctx, id)
	if err != nil {
		log.Printf("redisClient.GetOrder: %v", err)
		return nil, err
	}
	if order != nil {
		log.Printf("Found order with id [%s] in cache", id)
		return order, nil
	}

	if order, err = r.repo.GetById(ctx, id); err != nil {
		log.Printf("repo.GetById: %v", err)
		return nil, err
	}
	if order != nil {
		log.Printf("Found order with id [%s] in repo, now caching it", id)
		_ = r.redisClient.SaveOrder(ctx, order)
	}
	return order, nil
}
