package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TotalFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	From        time.Time
	To          time.Time
}

func (f *TotalFilter) Validate() error {
	if f.From.IsZero() || f.To.IsZero() {
		return errors.New("'from' and 'to' dates are required")
	}

	if f.To.Before(f.From) {
		return errors.New("'to' must be after 'from'")
	}

	return nil
}
