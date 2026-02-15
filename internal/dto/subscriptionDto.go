package dto

import "errors"

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

func (r *CreateSubscriptionRequest) Validate() error {
	if r.ServiceName == "" {
		return errors.New("service_name is required")
	}
	if r.Price <= 0 {
		return errors.New("price must be positive")
	}
	if r.UserID == "" {
		return errors.New("user_id is required")
	}
	if r.StartDate == "" {
		return errors.New("start_date is required")
	}
	return nil
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name"`
	Price       *int    `json:"price"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}
