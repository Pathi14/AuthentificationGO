package internal

import (
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)
var secretKey = []byte("my-very-secure-secret-key")

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")

        if !isValidToken(token) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func isValidToken(token string) bool {
    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })

    if err != nil {
        fmt.Println("Token parsing error:", err)
        return false
    }

    if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
        exp := int64(claims["exp"].(float64))
        if time.Now().Unix() > exp {
            fmt.Println("Token has expired")
            return false
        }
        return true
    }

    return false
}