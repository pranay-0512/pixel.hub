package services

import (
	"errors"
	"hub-api/db"
	"hub-api/models"
	"hub-api/utils"

	"gorm.io/gorm"
)

func CreateUser(name, username, password, email, role string) (*models.User, error) {
	var existingUser models.User

	if err := db.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err := db.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     name,
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Role:     role,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	user.Password = ""

	return &user, nil
}

func AuthenticateUser(username, password string) (string, error) {
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid password")
	}
	token := "dummy-token-for-" + user.Username
	return token, nil
}
