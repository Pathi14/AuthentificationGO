package api

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pathi14/AuthentificationGO/internal"
	"github.com/pathi14/AuthentificationGO/internal/infrastructure/database"
	"github.com/pathi14/AuthentificationGO/internal/middleware"
	"github.com/pathi14/AuthentificationGO/internal/user"
)

func Run() {
	router := gin.Default()

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	api := router.Group("/44df37e7-fe2a-404f-917b-399f5c5ffd12")
	{
		// Routes accessibles à tous
		api.GET("/health", internal.Health)
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.POST("/forgot-password", userHandler.ForgotPassword)
		api.POST("/reset-password", userHandler.ResetPassword)

		// Routes protégées
		api.Use(middleware.JWTAuth())
		{
			api.GET("/me", userHandler.Profile)
			api.POST("/logout", userHandler.Logout)
			api.POST("/refresh-token", userHandler.RefreshToken)
		}
	}

	fmt.Println("Server is listening on port 8080")

	port := os.Getenv("PORT") // Utilise le port fourni par Render

	if port == "" {
		port = "8080" // Port par défaut si non défini
	}

	router.Run("0.0.0.0:" + port) // Écoute sur 0.0.0.0 pour accepter les connexions externes
}
