package models

import "time"

type Map struct {
	ID          string       `json:"id" bson:"_id"`
	Name        string       `json:"name" bson:"name"`
	Width       int          `json:"width" bson:"width"`
	Height      int          `json:"height" bson:"height"`
	MapElements []MapElement `json:"map_elements" bson:"map_elements"`
	CreatedBy   string       `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" bson:"updated_at"`
}

type MapElement struct {
	ID   string `json:"id" bson:"_id"`
	PosX int    `json:"pos_x" bson:"pos_x"`
	PosY int    `json:"pos_y" bson:"pos_y"`
}
