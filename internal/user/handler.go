package user

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"message": "User added successfully",
		"user":    u,
	})
}

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "User logged in",
	})
}

func UpdatePassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "User password updated",
	})
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
    var request struct {
        Email string `json:"email" binding:"required,email"`
    }
 
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
 
    
    token, err := h.service.SendPasswordResetToken(request.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
 
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Password reset instructions sent to your email",
        "email":   request.Email,
        "token":   token, 
    })
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var request struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := h.service.ValidateResetToken(request.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	err = h.service.repo.ResetPassword(email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
		"NewPassword":   request.NewPassword,
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "User logged out",
	})
}

func Profile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "User profile",
	})
}
