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
		"message": "Utilisateur enregistré avec succès",
		"user":    u,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Une erreur interne est survenue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Connexion réussie",
		"token":   token,
	})
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Données d'entrée invalides",
			"details": err.Error(),
		})
		return
	}

	token, err := h.service.SendPasswordResetToken(request.Email)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusOK, gin.H{
				"message": "Si votre email est enregistré, vous recevrez un lien de réinitialisation.",
			})
			return
		}

		if strings.Contains(err.Error(), "validation error") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Une erreur est survenue lors de l'envoi de l'email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Les instructions pour la mise à jour de votre password ont été envoyés à votre email",
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

		if strings.Contains(err.Error(), "authentication error") ||
			strings.Contains(err.Error(), "invalid or expired token") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide ou expiré"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Une erreur système est survenue"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	if err := h.service.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Déconnexion échouée"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorisé"})
		return
	}

	user, err := h.service.GetUserByID(userID.(int))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Une erreur est survenue lors de la récupération du profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Profil de l'utilisateur",
		"email":       user.Email,
		"name":        user.Name,
		"age":         user.Age,
		"phoneNumber": user.MobileNumber,
	})
}
