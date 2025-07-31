package models

import (
	"time"

	"github.com/google/uuid"
)

type Sub struct {
	ID          uuid.UUID `json:"id"`
	UserID      int       `json:"user_id" validate:"required"`
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty"`
}

type SubRequest struct {
	ID          uuid.UUID `json:"id"`
	UserID      int       `json:"user_id" validate:"required"`
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty"`
}

type SumRequest struct {
	UserID      int    `json:"user_id"`
	ServiceName string `json:"service_name"`
	Start       string `json:"start"`
	End         string `json:"end"`
}
