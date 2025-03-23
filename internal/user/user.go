package user

import "github.com/go-playground/validator/v10"

type User struct {
	ID           int
	Name         string `json:"name"`
	Age          int    `json:"age" binding:"omitempty,gt=0"`
	MobileNumber string `json:"mobile_number"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password,omitempty"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
