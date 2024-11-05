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

type CreateMapReq struct {
	Name        string                                  `json:"name" bson:"name"`
	Width       int                                     `json:"width" bson:"width"`
	Height      int                                     `json:"height" bson:"height"`
	MapElements map[string]map[string]models.MapElement `json:"map_elements" bson:"map_elements"`
}

func CreateMap(c *gin.Context) {
	var req CreateMapReq
	var existingMap models.Map

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "user not Authenticated"})
	}

	maps := db.DB.Collection("maps")
	err := maps.FindOne(c, bson.M{"name": req.Name, "created_by": userId}).Decode(&existingMap)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Map name already taken"})
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
	id := uuid.New().String()

	roomMap := models.Map{
		ID:          id,
		Name:        req.Name,
		Width:       req.Width,
		Height:      req.Height,
		MapElements: req.MapElements,
		CreatedBy:   user.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err = maps.InsertOne(c, roomMap)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	_, err = db.DB.Collection("users").UpdateOne(c, bson.M{"_id": user.ID}, bson.M{"$push": bson.M{"maps": roomMap.ID}})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error, cannot push room to user"})
		return
	}

	utils.Success(c, gin.H{"map": roomMap})
}

func GetMaps(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "user not authenticated"})
		return
	}

	var user models.User
	err := db.DB.Collection("users").FindOne(c, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": "user not found"})
		return
	}

	var maps []models.Map
	cursor, err := db.DB.Collection("maps").Find(c, bson.M{"_id": bson.M{"$in": user.Maps}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err = cursor.All(c, &maps); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	utils.Success(c, gin.H{"maps": maps})
}

func UpdateMap(c *gin.Context) {

}

func GetMap(c *gin.Context) {
	mapId := c.Param("mapId")
	var roomMap models.Map

	_, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "user not authenticated"})
		return
	}
	err := db.DB.Collection("maps").FindOne(c, bson.M{"_id": mapId}).Decode(&roomMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	utils.Success(c, gin.H{"map": roomMap})
}

func GetAllMaps(c *gin.Context) {
	var maps []models.Map
	cursor, err := db.DB.Collection("maps").Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err = cursor.All(c, &maps); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	utils.Success(c, gin.H{"maps": maps})
}
