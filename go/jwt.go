/*
 * Grocery List Management API
 *
 * JWT token generation and validation
 */

package openapi

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Secret key for JWT signing - in production this should be an environment variable
var jwtSecretKey = []byte(getJWTSecret())

// getJWTSecret retrieves the JWT secret from environment or uses a default
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default secret for development - in production, always use environment variable
		return "dupa123dupa123dupa123dupa123dupa123dupa123"
	}
	return secret
}

// GenerateJWT generates a JWT token for the given user
func GenerateJWT(user *User) (string, error) {
	// Token expiration time (24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims with user information
	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "grocery-list-api",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT parses and validates a JWT token
func ParseJWT(tokenString string) (*JWTClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the token and extract claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// JWTAuthMiddleware is a middleware for JWT authentication
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, Error{
				Message: "Authorization header missing",
				Details: "Bearer token required",
			})
			c.Abort()
			return
		}

		// The token should be in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, Error{
				Message: "Invalid authorization format",
				Details: "Authorization header must be in the format: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Parse and validate the token
		tokenString := parts[1]
		claims, err := ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, Error{
				Message: "Invalid or expired token",
				Details: err.Error(),
			})
			c.Abort()
			return
		}

		// Set the user ID and username in the context for use in handlers
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
