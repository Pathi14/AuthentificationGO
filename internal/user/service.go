package user

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/dgrijalva/jwt-go"
	"github.com/pathi14/AuthentificationGO/internal/infrastructure/database"
	"github.com/pathi14/AuthentificationGO/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo        *UserRepository
	tokenExpiry time.Duration
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{
		repo:        repo,
		tokenExpiry: 15 * time.Minute,
	}
}

func (s *UserService) Create(u User) error {
	if err := u.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	existingUser, err := s.repo.GetByEmail(u.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("internal error: %v", err)
	}
	if existingUser != nil {
		return fmt.Errorf("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	u.Password = string(hashedPassword)

	if err := s.repo.Create(u); err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return fmt.Errorf("email already in use")
		}
		return fmt.Errorf("internal error: %v", err)
	}
	return nil
}

func (s *UserService) Login(email, password string) (string, string, error) {
	if email == "" {
		return "", "", fmt.Errorf("validation error: email is required")
	}
	if password == "" {
		return "", "", fmt.Errorf("validation error: password is required")
	}

	user, err := s.repo.Login(email, password)
	if err != nil {
		if err == sql.ErrNoRows || strings.Contains(err.Error(), "user not found") {
			return "", "", fmt.Errorf("authentication error: user not found")
		}
		if strings.Contains(err.Error(), "invalid password") || strings.Contains(err.Error(), "password does not match") {
			return "", "", fmt.Errorf("authentication error: invalid credentials")
		}
		return "", "", fmt.Errorf("internal error: %v", err)
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", "", fmt.Errorf("internal error: JWT_SECRET not configured")
	}

	accessToken, err := s.generateToken(int(user.ID), 2*time.Hour)
	if err != nil {
		return "", "", fmt.Errorf("internal error: failed to generate access token")
	}

	newRefreshToken, err := s.generateToken(int(user.ID), 7*24*time.Hour)
	if err != nil {
		return "", "", fmt.Errorf("internal error: failed to generate refresh token")
	}

	return accessToken, newRefreshToken, nil
}

func (s *UserService) Logout(tokenString string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return fmt.Errorf("erreur de connexion à la base de données : %v", err)
	}
	defer db.Close()

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("malformed token")
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

	return middleware.AddToBlacklist(db, tokenString, expirationTime)
}

func (s *UserService) ResetPassword(token, newPassword string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return fmt.Errorf("erreur de connexion à la base de données : %v", err)
	}
	defer db.Close()

	if token == "" {
		return fmt.Errorf("validation error: token is required")
	}
	if newPassword == "" {
		return fmt.Errorf("validation error: new password is required")
	}
	if len(newPassword) < 8 {
		return fmt.Errorf("validation error: password must be at least 8 characters long")
	}

	if middleware.IsTokenBlacklisted(token) {
		return fmt.Errorf("authentication error: token already used or expired")
	}

	email, err := s.ValidateResetToken(token)
	if err != nil {
		return fmt.Errorf("authentication error: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("internal error: failed to hash password: %v", err)
	}

	err = s.repo.ResetPassword(email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("internal error: failed to update password: %v", err)
	}

	expirationTime := time.Now().Add(s.tokenExpiry)
	return middleware.AddToBlacklist(db, token, expirationTime)
}

func (s *UserService) SendPasswordResetToken(email string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("validation error: email is required")
	}

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("internal error: %v", err)
	}

	if user == nil {
		return "", fmt.Errorf("user not found")
	}

	resetToken, err := generateResetToken(user.Email)
	if err != nil {
		return "", fmt.Errorf("internal error: failed to generate reset token: %v", err)
	}

	err = sendResetEmail(user.Email, resetToken)
	if err != nil {
		return "", fmt.Errorf("internal error: failed to send reset email: %v", err)
	}

	return resetToken, nil
}

func generateResetToken(email string) (string, error) {
	if email == "" {
		return "", errors.New("email cannot be empty")
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}

func (s *UserService) ValidateResetToken(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("authentication error: token is required")
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("authentication error: invalid token: %v", err)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if email, ok := claims["email"].(string); ok {
			return email, nil
		}
		return "", fmt.Errorf("authentication error: invalid token structure")
	}
	return "", fmt.Errorf("authentication error: invalid token")
}

func sendResetEmail(email, token string) error {
	if email == "" || token == "" {
		return errors.New("email and token are required")
	}

	fmt.Printf("To: %s\n", email)
	fmt.Printf("Subject: Password Reset Request\n")
	fmt.Printf("Body: Click the link below to reset your password:\n")
	fmt.Printf("https://go-auth-api-latest.onrender.com/44df37e7-fe2a-404f-917b-399f5c5ffd12/reset-password?token=%s\n", token)
	return nil
}

func (s *UserService) GetUserByID(id int) (*User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("validation error: invalid user ID")
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("not found: user with ID %d does not exist", id)
		}
		return nil, fmt.Errorf("internal error: %v", err)
	}

	if user == nil {
		return nil, fmt.Errorf("not found: user with ID %d does not exist", id)
	}

	return user, nil
}

func (s *UserService) refreshToken(refreshToken string) (string, string, error) {

	if refreshToken == "" {
		return "", "", fmt.Errorf("validation error: refresh token is required")
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", "", fmt.Errorf("internal error: JWT_SECRET not configured")
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", "", fmt.Errorf("authentication error: invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("authentication error: invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return "", "", fmt.Errorf("authentication error: invalid user_id")
	}

	expiration, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(expiration) {
		return "", "", fmt.Errorf("authentication error: refresh token expired")
	}

	accessToken, err := s.generateToken(int(userID), 2*time.Hour)
	if err != nil {
		return "", "", fmt.Errorf("internal error: failed to generate access token")
	}

	newRefreshToken, err := s.generateToken(int(userID), 7*24*time.Hour)
	if err != nil {
		return "", "", fmt.Errorf("internal error: failed to generate refresh token")
	}

	db, err := database.ConnectDB()
	if err != nil {
		return "", "", fmt.Errorf("erreur de connexion à la base de données : %v", err)
	}
	defer db.Close()

	expirationTime := time.Now().Add(s.tokenExpiry)
	middleware.AddToBlacklist(db, refreshToken, expirationTime)

	return accessToken, newRefreshToken, nil
}

func (s *UserService) generateToken(userID int, duration time.Duration) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("internal error: JWT_SECRET not configured")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"id":      uuid.New().String(),
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
