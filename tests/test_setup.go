package tests

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pathi14/AuthentificationGO/internal/infrastructure/database"
	"github.com/pathi14/AuthentificationGO/internal/user"
)

var testRouter *gin.Engine

//var db *database.DB

func init() {
	os.Setenv("TEST_DB_HOST", "localhost")
	os.Setenv("TEST_DB_PORT", "5434")
	os.Setenv("TEST_DB_USER", "test_admin")
	os.Setenv("TEST_DB_PASSWORD", "test_secret")
	os.Setenv("TEST_DB_NAME", "authentificationgo_test")

	db, err := database.ConnectTestDB()
	if err != nil {
		panic("Erreur lors de la connexion à la DB de test : " + err.Error())
	}

	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	testRouter = gin.Default()
	testRouter.POST("/login", userHandler.Login)
	testRouter.POST("/register", userHandler.Register)
}

// func TestMain(m *testing.M) {
// 	code := m.Run()
// 	db.Exec("DELETE FROM users")
// 	db.Close()
// 	os.Exit(code)
// }
