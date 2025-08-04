package dto

type UserRequest struct {
	Username  string `json:"username" validate:"required,min=3"`
	Password  string `json:"password" validate:"required,min=6"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}
