/*
 * Grocery List Management API
 *
 * Authentication API handlers
 */

package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthAPI struct {
}

// Post /register
// Register a new user
func (api *AuthAPI) Register(c *gin.Context) {
	var userRegister UserRegister

	// Bind the request body to the UserRegister struct
	if err := c.ShouldBindJSON(&userRegister); err != nil {
		c.JSON(http.StatusBadRequest, Error{
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Register the user
	user, err := RegisterUser(userRegister)
	if err != nil {
		// Determine the appropriate error status code
		statusCode := http.StatusInternalServerError
		if err.Error() == "username already exists" ||
			err.Error() == "username is required" ||
			err.Error() == "password is required" ||
			err.Error() == "repeated_password is required" ||
			err.Error() == "passwords do not match" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, Error{
			Message: "Registration failed",
			Details: err.Error(),
		})
		return
	}

	// Return the user without password
	c.JSON(http.StatusCreated, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"createdAt": user.CreatedAt,
	})
}

// Post /login
// Login a user
func (api *AuthAPI) Login(c *gin.Context) {
	var userLogin UserLogin

	// Bind the request body to the UserLogin struct
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, Error{
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Login the user
	user, err := LoginUser(userLogin)
	if err != nil {
		// For security reasons, we don't want to reveal too much information about the error
		c.JSON(http.StatusUnauthorized, Error{
			Message: "Authentication failed",
			Details: "Invalid username or password",
		})
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{
			Message: "Token generation failed",
			Details: "Could not generate authentication token",
		})
		return
	}

	// Return the user and token
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"createdAt": user.CreatedAt,
		},
		"token": token,
	})
}
