package user

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	 "github.com/dgrijalva/jwt-go"
 
)
var secretKey = []byte("my-very-secure-secret-key")
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

func (s *UserService) Login(email, password string) (string, error) {
    user, err := s.repo.Login(email, password)
    if err != nil {
        return "", err
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString(secretKey) // Utilisez la même clé secrète
    if err != nil {
        return "", fmt.Errorf("erreur lors de la génération du token: %v", err)
    }

    return tokenString, nil
}