package services

import (
	"errors"
	"fmt"

	"github.com/SadikMR/go-expense-tracker-api/models"
)

var (
	ErrEmailExists  = errors.New("email already exists")
	ErrInvalidCreds = errors.New("invalid email or password")
)

// Register creates a new user after validating uniqueness.
// Returns ErrEmailExists if the email is already taken.
func Register(name, email, password string) error {
	existing, err := models.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("Register: %w", err)
	}
	if existing != nil {
		return ErrEmailExists
	}

	id, err := models.NextUserID()
	if err != nil {
		return fmt.Errorf("Register: %w", err)
	}

	if err := models.CreateUser(&models.User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}); err != nil {
		return fmt.Errorf("Register: %w", err)
	}
	return nil
}

// Login verifies credentials and returns the authenticated user.
// Returns ErrInvalidCreds if email does not exist or password does not match.
func Login(email, password string) (*models.User, error) {
	user, err := models.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}
	if user == nil || user.Password != password {
		return nil, ErrInvalidCreds
	}
	return user, nil
}
