package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *UserRepository
}

type ResetToken struct {
	Token     string
	Email     string
	ExpiresAt time.Time
}

var resetTokens = make(map[string]ResetToken)

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


func (s *UserService) ForgotPassword(email string) error {
	user, err := s.repo.GetByEmail(email)
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}

	// Génération d'un token unique
	token := uuid.New().String()
	resetTokens[token] = ResetToken{
		Token:     token,
		Email:     email,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Expiration en 1 heure
	}

	// Simuler l'envoi d'un email (dans la vraie vie, envoyer par email)
	fmt.Printf("Reset token for %s: %s\n", email, token)

	return nil
}

func (s *UserService) ResetPassword(token, newPassword string) error {
	resetData, exists := resetTokens[token]
	if !exists || time.Now().After(resetData.ExpiresAt) {
		return fmt.Errorf("invalid or expired token")
	}

	// Hash du nouveau mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	// Mise à jour du mot de passe dans la BD
	err = s.repo.UpdatePassword(resetData.Email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("error updating password: %v", err)
	}

	// Suppression du token après usage
	delete(resetTokens, token)
	return nil
}