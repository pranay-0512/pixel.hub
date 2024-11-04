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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateRoomReq struct {
	Name        string `json:"name" binding:"required" bson:"name"`
	Description string `json:"description" bson:"description"`
	Thumbnail   string `json:"thumbnail" bson:"thumbnail"`
	RoomMap     string `json:"room_map" bson:"room_map"`
	Capacity    int    `json:"capacity" bson:"capacity"`
	IsPrivate   bool   `json:"is_private" bson:"is_private"`
}

func CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	var existingRoom models.Room

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "user not authenticated"})
		return
	}
	rooms := db.DB.Collection("rooms")
	err := rooms.FindOne(c, bson.M{"name": req.Name}).Decode(&existingRoom)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Room name already taken"})
		return
	}

	var user models.User
	err = db.DB.Collection("users").FindOne(c, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": "user not found"})
		return
	}

	id := uuid.New().String()
	room := models.Room{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
		RoomMap:     req.RoomMap,
		Capacity:    req.Capacity,
		IsPrivate:   req.IsPrivate,
		Owner:       user.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = rooms.InsertOne(c, room)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	_, err = db.DB.Collection("users").UpdateOne(c, bson.M{"_id": user.ID}, bson.M{"$push": bson.M{"rooms": room.ID}})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error, cannot push room to user"})
		return
	}

	utils.Success(c, gin.H{"room": room})
}

func GetRooms(c *gin.Context) {
	// get user id from context
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

	// return all the rooms from user.Rooms
	var rooms []models.Room
	cursor, err := db.DB.Collection("rooms").Find(c, bson.M{"_id": bson.M{"$in": user.Rooms}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err = cursor.All(c, &rooms); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	utils.Success(c, gin.H{"rooms": rooms})
}

func GetRoom(c *gin.Context) {
	roomId := c.Param("roomId")
	var room models.Room

	_, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "user not authenticated"})
		return
	}

	err := db.DB.Collection("rooms").FindOne(c, bson.M{"_id": roomId}).Decode(&room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	utils.Success(c, gin.H{"room": room})
}

func UpdateRoom(c *gin.Context) {
	roomId := c.Param("roomId")
	var req CreateRoomReq
	var existingRoom models.Room

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "user not authenticated"})
		return
	}

	after := options.After
	err := db.DB.Collection("rooms").FindOneAndUpdate(c, bson.M{"_id": roomId}, bson.M{"$set": bson.M{"name": req.Name, "description": req.Description, "thumbnail": req.Thumbnail, "room_map": req.RoomMap, "capacity": req.Capacity, "is_private": req.IsPrivate, "updated_at": time.Now()}}, &options.FindOneAndUpdateOptions{ReturnDocument: &after}).Decode(&existingRoom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	utils.Success(c, gin.H{"room": existingRoom})
}

func GetAllRooms(c *gin.Context) {
	var rooms []models.Room
	cursor, err := db.DB.Collection("rooms").Find(c, bson.M{"is_private": false})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err = cursor.All(c, &rooms); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	utils.Success(c, gin.H{"rooms": rooms})
}
