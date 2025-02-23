package helpers

import (
	"backend/graph/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func FindUserByEmail(email string) (*model.User, error) {
	// Mock database lookup
	if email == "test@example.com" {
		return &model.User{
			ID:       "1",
			Username: "Test User",
			Email:    "test@example.com",
			//	Password: "$2a$10$ExampleHashedPassword", // Replace with actual hashed password
		}, nil
	}
	return nil, errors.New("user not found")
}

func EmailExists(email string) bool {
	// Mock email existence check
	return email == "test@example.com"
}

func SaveUser(user *model.User) error {
	// Mock save user logic
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
