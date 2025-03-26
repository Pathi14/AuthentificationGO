package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pathi14/AuthentificationGO/internal/infrastructure/database"
)

// JWTAuth vérifie le token JWT et extrait l'ID utilisateur
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer le token depuis l'en-tête Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort() // Important : arrête l'exécution des middlewares suivants
			return
		}

		// Format Bearer token
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if IsTokenBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalid or expired"})
			c.Abort()
			return
		}

		// Valider et décoder le token
		secretKey := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Vérifier l'algorithme de signature
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Retourner la clé secrète pour la vérification
			return []byte(secretKey), nil // Stockez cette clé dans une variable d'environnement
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// Extraire l'ID utilisateur du token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp, ok := claims["exp"].(float64)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token expiration"})
				c.Abort()
				return
			}
			// Vérifie si le token est expiré
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				c.Abort()
				return
			}

			userID := int(userIDFloat)
			c.Set("userID", userID) // Stocke l'ID utilisateur dans le contexte
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

func AddToBlacklist(db *sql.DB, token string, expiration time.Time) error {
	_, err := db.Exec("INSERT INTO blacklisted_tokens (token, expiration) VALUES ($1, $2)", token, expiration)
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}
	return nil
}

func IsTokenBlacklisted(token string) bool {
	db, err := database.ConnectDB()
	if err != nil {
		return true
	}
	defer db.Close()

	var expiration time.Time
	err = db.QueryRow("SELECT expiration FROM blacklisted_tokens WHERE token = $1", token).Scan(&expiration)

	if err != nil {
		return false
	}

	if time.Now().After(expiration) {
		_, err := db.Exec("DELETE FROM blacklisted_tokens WHERE token = $1", token)
		if err != nil {
			return true
		}
		return false
	}

	return true
}
