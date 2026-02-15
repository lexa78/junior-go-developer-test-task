package repository

import (
	"context"
	"testTask/internal/domain"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, s *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, userId *uuid.UUID, serviceName *string) ([]domain.Subscription, error)
	CalculateTotal(ctx context.Context, filter *domain.TotalFilter) (int, error)
}
