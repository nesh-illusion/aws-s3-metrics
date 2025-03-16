package services

import (
	"project_security_one/internal/models"
	"project_security_one/internal/repositories"
)

// UserService handles business logic
type UserService struct {
	UserRepo *repositories.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(user *models.User) error {
	return s.UserRepo.CreateUser(user)
}
