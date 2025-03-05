package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pathi14/AuthentificationGO/internal"
)

func Run() {
	router := gin.Default()
	router.GET("/health", internal.Health)

	router.Run("localhost:8080")
}
