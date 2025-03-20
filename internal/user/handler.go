package user

import (
	"net/http"

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
