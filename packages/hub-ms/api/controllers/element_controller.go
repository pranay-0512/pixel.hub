package controllers

import (
	"hub-api/db"
	"hub-api/models"
	"hub-api/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateElementReq struct {
	Name    string `json:"name" binding:"required" bson:"name"`
	Enum    string `json:"enum" binding:"required" bson:"enum"`
	Width   int    `json:"width" binding:"required" bson:"width"`
	Height  int    `json:"height" binding:"required" bson:"height"`
	URL     string `json:"url" binding:"required" bson:"url"`
	IsSolid bool   `json:"is_solid" binding:"required" bson:"is_solid"`
}

func CreateElement(c *gin.Context) {
	var req CreateElementReq
	var existingEle models.Element

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "user not Authenticated"})
		return
	}

	elements := db.DB.Collection("elements")
	err := elements.FindOne(c, bson.M{"name": req.Name}).Decode(&existingEle)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "element with this name exists"})
		return
	}
	var user models.User
	err = db.DB.Collection("users").FindOne(c, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": "user not found"})
		return
	}
	if user.Role != "ADMIN" {
		c.JSON(400, gin.H{"error": "Not Authorised to create a map"})
		return
	}

	element := models.Element{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Enum:      req.Enum,
		Width:     req.Width,
		Height:    req.Height,
		URL:       req.URL,
		IsSolid:   req.IsSolid,
		CreatedBy: user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = elements.InsertOne(c, element)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error, cannot insert element"})
		return
	}
	_, err = db.DB.Collection("users").UpdateOne(c, bson.M{"_id": user.ID}, bson.M{"$push": bson.M{"maps": element.ID}})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error, cannot push room to user"})
		return
	}

	utils.Success(c, gin.H{"element": element})
}

func GetElements(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "user not Authenticated"})
		return
	}

	var user models.User
	err := db.DB.Collection("users").FindOne(c, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": "user not found"})
		return
	}

	var elements []models.Element
	cursor, err := db.DB.Collection("elements").Find(c, bson.M{"_id": bson.M{"$in": user.Elements}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err = cursor.All(c, &elements); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	utils.Success(c, gin.H{"elements": elements})
}

func GetElement(c *gin.Context) {
	elementId := c.Param("id")
	var element models.Element

	_, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "user not Authenticated"})
		return
	}
	err := db.DB.Collection("elements").FindOne(c, bson.M{"_id": elementId}).Decode(&element)
	if err != nil {
		c.JSON(400, gin.H{"error": "element not found"})
		return
	}

	utils.Success(c, gin.H{"element": element})
}

func GetAllElements(c *gin.Context) {
	var elements []models.Element
	cursor, err := db.DB.Collection("elements").Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err = cursor.All(c, &elements); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	utils.Success(c, gin.H{"elements": elements})
}

func UpdateElement(c *gin.Context) {}
