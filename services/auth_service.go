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
	"gorm.io/gorm"
)

func Register(user *models.User) (*models.User, int, error) {
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return loginAndGenerateToken(&existingUser, user.Password)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	if err := database.DB.Create(user).Error; err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create user: %v", err)
	}

	return loginAndGenerateToken(user, user.Password)
}

func Login(email string, password string) (*models.User, int, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, http.StatusUnauthorized, errors.New("invalid email or password")
		}
		return nil, http.StatusInternalServerError, errors.New("database error")
	}

	return loginAndGenerateToken(&user, password)
}

// loginAndGenerateToken handles user login and token generation.
func loginAndGenerateToken(user *models.User, password string) (*models.User, int, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("invalid email or password")
	}

	tokenString, err := utility.GenerateToken(user.ID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to create token")
	}

	user.Token = tokenString
	return user, http.StatusOK, nil
}

func RefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil && err != jwt.ErrTokenExpired {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrInvalidKey
	}

	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
