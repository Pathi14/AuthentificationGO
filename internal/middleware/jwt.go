// internal/middleware/jwt.go
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4" // Importez cette bibliothèque
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

		// Valider et décoder le token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Vérifier l'algorithme de signature
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Retourner la clé secrète pour la vérification
			return []byte("votre_clé_secrète"), nil // Stockez cette clé dans une variable d'environnement
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// Extraire l'ID utilisateur du token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Convertir l'ID en int
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				c.Abort()
				return
			}

			userID := int(userIDFloat)
			c.Set("userID", userID) // Stocke l'ID utilisateur dans le contexte
			c.Next()                // Continue vers le handler suivant
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
