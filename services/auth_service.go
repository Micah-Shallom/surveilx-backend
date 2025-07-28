package services

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"survielx-backend/database"
	"survielx-backend/models"
	"survielx-backend/utility"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(user *models.User) (int,error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusBadRequest,err
	}
	user.ID = utility.GenerateUUID()

	user.Password = string(hashedPassword)
	result := database.DB.Create(user)
	if result.Error != nil {
		return http.StatusInternalServerError, fmt.Errorf("Failed to create user")
	}
	return http.StatusOK, nil
}

func Login(email string, password string) (string, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return "", errors.New("invalid email or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
