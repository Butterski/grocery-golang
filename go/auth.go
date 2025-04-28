/*
 * Grocery List Management API
 *
 * User model and authentication functions
 */

package openapi

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the database
type User struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"` // Password is not exported to JSON
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// UserRegister represents the registration request
type UserRegister struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	RepeatedPassword string `json:"repeated_password"`
}

// UserLogin represents the login request
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate validates the registration request
func (r *UserRegister) Validate() error {
	if r.Username == "" {
		return errors.New("username is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if r.RepeatedPassword == "" {
		return errors.New("repeated_password is required")
	}
	if r.Password != r.RepeatedPassword {
		return errors.New("passwords do not match")
	}
	return nil
}

// Validate validates the login request
func (l *UserLogin) Validate() error {
	if l.Username == "" {
		return errors.New("username is required")
	}
	if l.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// RegisterUser registers a new user
func RegisterUser(userReg UserRegister) (*User, error) {
	// Validate the registration request
	if err := userReg.Validate(); err != nil {
		return nil, err
	}

	// Check if user already exists
	var existingUser User
	result := DB.Where("username = ?", userReg.Username).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("username already exists")
	} else if result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReg.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create a new user
	user := User{
		Username: userReg.Username,
		Password: string(hashedPassword),
	}

	// Save the user to the database
	result = DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// LoginUser authenticates a user
func LoginUser(userLogin UserLogin) (*User, error) {
	// Validate the login request
	if err := userLogin.Validate(); err != nil {
		return nil, err
	}

	// Find the user
	var user User
	result := DB.Where("username = ?", userLogin.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid username or password")
		}
		return nil, result.Error
	}

	// Check the password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}
