package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(u User) error {
	if u.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if u.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if u.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	existingUser, err := s.repo.GetByEmail(u.Email)
	if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}
	if existingUser != nil {
		return fmt.Errorf("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	u.Password = string(hashedPassword)

	return s.repo.Create(u)
}
