package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pathi14/AuthentificationGO/internal"
	"github.com/pathi14/AuthentificationGO/internal/infrastructure/database"
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
		api.POST("/login", user.Login)
		api.POST("/reset-password", user.UpdatePassword)

		// Routes protégées
		api.POST("/logout", user.Logout)
		api.GET("/profile", user.Profile)
	}

	fmt.Println("Server is listening on port 8080")
	router.Run("localhost:8080")
}
