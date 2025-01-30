package usecases

import (
	"errors"
	"time"

	"main.go/entities"
	"main.go/interface_adapters/repositories"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type LoginUser struct {
	UserRepo repositories.UserRepository
	JWTKey   []byte
}

func (l *LoginUser) Execute(user entities.User) (string, error) {
	storedUser, exists := l.UserRepo.FindByUsername(user.Username)
	if !exists {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &entities.User.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(l.JWTKey)
	if err != nil {
		return "", errors.New("could not generate token")
	}

	return tokenString, nil
}
