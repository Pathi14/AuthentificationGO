package user

import "github.com/go-playground/validator/v10"

type User struct {
	ID           int
	Name         string `json:"name" binding:"required,min=2,max=50"`
	Age          int    `json:"age" binding:"omitempty,gt=0"`
	MobileNumber string `json:"mobile_number"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password,omitempty" binding:"required,min=8"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
