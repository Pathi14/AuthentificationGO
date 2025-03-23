package user
 
import (
    "fmt"
    "time"
    "sync"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
 
)
 
type UserService struct {
    repo           *UserRepository
    resetTokens    map[string]time.Time 
    tokenMutex     sync.Mutex           
    tokenExpiry    time.Duration        
}
 
func NewUserService(repo *UserRepository) *UserService {
    return &UserService{
        repo:        repo,
        resetTokens: make(map[string]time.Time),
        tokenExpiry: 15 * time.Minute, 
    }
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
 
    tokenString, err := token.SignedString([]byte("secret_key"))
    if err != nil {
        return "", fmt.Errorf("erreur lors de la génération du token: %v", err)
    }
 
    return tokenString, nil
}
 
func (s *UserService) SendPasswordResetToken(email string) (string, error) {
   
    user, err := s.repo.GetByEmail(email)
    if err != nil {
        return "", errors.New("user not found")
    }
 
    
    resetToken, err := generateResetToken(user.Email)
    if err != nil {
        return "", fmt.Errorf("failed to generate reset token: %v", err)
    }
 
    
    err = sendResetEmail(user.Email, resetToken)
    if err != nil {
        return "", errors.New("failed to send reset email")
    }
 
    return resetToken, nil
}
 
func generateResetToken(email string) (string, error) {
    
    secretKey := []byte("secret_key")
 
    
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
   
    secretKey := []byte("secret_key")
 
    
    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })
 
    if err != nil {
        return "", fmt.Errorf("invalid token: %v", err)
    }
 
    
    if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
        email := claims["email"].(string) 
        return email, nil
    }
 
    return "", errors.New("invalid token")
}
 
func sendResetEmail(email, token string) error {
    
    fmt.Printf("To: %s\n", email)
    fmt.Printf("Subject: Password Reset Request\n")
    fmt.Printf("Body: Click the link below to reset your password:\n")
    fmt.Printf("https://your-app.com/reset-password?token=%s\n", token)
    return nil
}


func (s *UserService) ResetPassword(token, newPassword string) error {
    email, err := s.ValidateResetToken(token)
    if err != nil {
        return errors.New("invalid or expired token")
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return errors.New("error hashing password")
    }

    err = s.repo.ResetPassword(email, string(hashedPassword))
    if err != nil {
        return errors.New("error password")
    }

    return nil
}


