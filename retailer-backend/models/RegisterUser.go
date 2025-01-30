package usecases

import (
	"errors"
	"retailer-backend/entities"
	"retailer-backend/interface_adapters/repositories"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	UserRepo repositories.UserRepository
}

func (r *RegisterUser) Execute(user entities.User) error {
	// Check if the user already exists
	if _, exists := r.UserRepo.FindByUsername(user.Username); exists {
		return errors.New("username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("could not hash password")
	}

	// Store the user
	user.Password = string(hashedPassword)
	r.UserRepo.Save(user)
	return nil
}
