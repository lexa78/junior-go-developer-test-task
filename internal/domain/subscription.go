package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

func (s *Subscription) Validate() error {
	if s.ServiceName == "" {
		return errors.New("service name is required")
	}

	if s.Price <= 0 {
		return errors.New("price must be positive")
	}

	if s.StartDate.IsZero() {
		return errors.New("start date is required")
	}

	if s.EndDate != nil && s.EndDate.Before(s.StartDate) {
		return errors.New("end date cannot be before start date")
	}

	return nil
}
