package user

import (
	"net/http"
	"strings"

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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Données d'entrée invalides",
			"details": err.Error(),
		})
		return
	}

	err := h.service.Create(u)
	if err != nil {
		if strings.Contains(err.Error(), "validation error") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if strings.Contains(err.Error(), "email already in use") {
			c.JSON(http.StatusConflict, gin.H{"error": "Cet email est déjà utilisé"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"message": "User added successfully",
		"user":    u,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Données d'entrée invalides",
			"details": err.Error(),
		})
		return
	}

	token, err := h.service.Login(credentials.Email, credentials.Password)
	if err != nil {

		if strings.Contains(err.Error(), "validation error") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if strings.Contains(err.Error(), "authentication error") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Une erreur interne est survenue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"token":   token,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Données d'entrée invalides",
			"details": err.Error(),
		})
		return
	}

	err := h.service.ResetPassword(request.Token, request.NewPassword)
	if err != nil {
		if strings.Contains(err.Error(), "validation error") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if strings.Contains(err.Error(), "authentication error") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Une erreur système est survenue"})
		return
	}

	c.Status(http.StatusNoContent)
}

func Logout(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.service.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message":     "User profile",
		"email":       user.Email,
		"name":        user.Name,
		"lastName":    user.Age,
		"phoneNumber": user.MobileNumber,
	})
}
