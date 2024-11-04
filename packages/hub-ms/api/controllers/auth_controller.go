package controllers

import (
	"fmt"
	"hub-api/config"
	"hub-api/db"
	"hub-api/models"
	"hub-api/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignUpReq struct {
	Username string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role" bson:"role"`
}

type LoginReq struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

var jwtSecret = []byte(config.AppConfig.JWTSecret)

func SignUp(c *gin.Context) {
	var req SignUpReq
	var existingUser *models.User
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	users := db.DB.Collection("users")
	err := users.FindOne(ctx, bson.M{"username": req.Username}).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Name already taken"})
		return
	}

	hashedPass, err := utils.HashPassword(req.Password)
	if err != nil {
		panic(err)
	}
	id := uuid.New().String()

	if req.Role == "ADMIN" {
		user := models.User{
			ID:        id,
			Name:      req.Name,
			Username:  req.Username,
			Email:     req.Email,
			Password:  hashedPass,
			Role:      req.Role,
			Elements:  []string{},
			Maps:      []string{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err = users.InsertOne(ctx, user)
		if err != nil {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		user.Password = ""
		utils.Success(c, gin.H{"user": user})
	} else {
		user := models.User{
			ID:        id,
			Name:      req.Name,
			Username:  req.Username,
			Email:     req.Email,
			Password:  hashedPass,
			Role:      req.Role,
			Rooms:     []string{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err = users.InsertOne(ctx, user)
		if err != nil {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		user.Password = ""
		utils.Success(c, gin.H{"user": user})
	}

}

func Login(c *gin.Context) {
	var req LoginReq
	var existingUser *models.User
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	users := db.DB.Collection("users")
	err := users.FindOne(ctx, bson.M{"username": req.Username}).Decode(&existingUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if !utils.CheckPasswordHash(req.Password, existingUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := GenerateToken(existingUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	utils.Success(c, gin.H{"token": token})
}

func GenerateToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userId,
		"exp":    time.Now().Add(time.Hour * 48).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
