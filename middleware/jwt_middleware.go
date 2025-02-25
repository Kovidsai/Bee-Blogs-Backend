package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT Secret Key
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// JWT Middleware to protect routes
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Extract token
		tokenStringSlice := strings.Split(authHeader, "Bearer ")
		if len(tokenStringSlice) == 1 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No JWT"})
			c.Abort()
			return
		}
		tokenString := tokenStringSlice[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Get `user_id` from token claims
		userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID missing in token"})
			c.Abort()
			return
		}

		// Store `userID` in Gin context
		c.Set("userID", int(userID)) // Convert float64 to int
		c.Next()
	}
}
