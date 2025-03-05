package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pathi14/AuthentificationGO/internal"
	"github.com/pathi14/AuthentificationGO/internal/user"
)

func Run() {
	router := gin.Default()

	//routes accessibles à tous
	router.GET("/health", internal.Health)
	router.POST("/register", user.Register)
	router.POST("/login", user.Login)
	router.POST("/reset-password", user.UpdatePassword)

	//routes protégées
	router.POST("/logout", user.Logout)
	router.GET("/profile", user.Profile)

	router.Run("localhost:8080")
}
