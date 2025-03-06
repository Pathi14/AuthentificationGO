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

	//routes accessibles à tous
	router.GET("/health", internal.Health)
	router.POST("/register", userHandler.Register)
	router.POST("/login", user.Login)
	router.POST("/reset-password", user.UpdatePassword)

	//routes protégées
	router.POST("/logout", user.Logout)
	router.GET("/profile", user.Profile)

	fmt.Println("Server is listening on port 8080")
	router.Run("localhost:8080")
}
