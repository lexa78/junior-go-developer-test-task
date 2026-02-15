package service

import (
	"context"
	"testTask/internal/domain"
	"testTask/internal/repository"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(ctx context.Context, sub *domain.Subscription) error {
	sub.ID = uuid.New()

	if err := sub.Validate(); err != nil {
		return err
	}

	return s.repo.Create(ctx, sub)
}

func (s *SubscriptionService) List(ctx context.Context, userId *uuid.UUID, serviceName *string) ([]domain.Subscription, error) {
	return s.repo.List(ctx, userId, serviceName)
}

func (s *SubscriptionService) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubscriptionService) Update(ctx context.Context, sub *domain.Subscription) error {
	if err := sub.Validate(); err != nil {
		return err
	}

	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *SubscriptionService) CalculateTotal(
	ctx context.Context,
	filter *domain.TotalFilter,
) (int, error) {
	return s.repo.CalculateTotal(ctx, filter)
}
