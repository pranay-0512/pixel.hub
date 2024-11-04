package models

import "time"

type Room struct {
	ID          string    `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Thumbnail   string    `json:"thumbnail" bson:"thumbnail"`
	RoomMap     string    `json:"room_map" bson:"room_map"`
	Capacity    int       `json:"capacity" bson:"capacity"`
	IsPrivate   bool      `json:"is_private" bson:"is_private"`
	Owner       string    `json:"owner" bson:"owner"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
