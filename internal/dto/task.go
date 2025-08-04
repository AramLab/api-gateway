package dto

import "time"

type TaskResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskRequest struct {
	UserID      string `json:"user_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"oneof=new in_progress done"`
}
