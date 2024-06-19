package utils

type RegisterInput struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Companyname string `json:"companyName" binding:"required"`
	Password    string `json:"password" binding:"required,min=8" validate:"password_complexity"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8" validate:"password_complexity"`
}
