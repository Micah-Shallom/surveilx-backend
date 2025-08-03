package services

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"survielx-backend/database"
	"survielx-backend/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, user *models.User) (*models.User, int, error) {
	var (
		existingUser models.User
		profile      models.Profile
	)

	originalPassword := user.Password

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to start transaction")
	}

	exist := models.CheckExists(database.DB, &existingUser, "email = ?", user.Email)
	if exist {
		return nil, http.StatusConflict, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	err = user.CreateUser(tx)
	if err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create user: %v", err)
	}

	profile.UserID = user.ID
	err = profile.CreateProfile(tx)
	if err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create user profile: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to commit transaction")
	}

	luser, code, err := loginAndGenerateToken(user, originalPassword)
	if err != nil {
		return nil, code, fmt.Errorf("failed to login and generate token: %v", err)
	}


	return luser, http.StatusCreated, nil
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
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
