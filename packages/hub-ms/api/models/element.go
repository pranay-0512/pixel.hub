package models

import "time"

type Element struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Enum      string    `json:"enum" bson:"enum"`
	Width     int       `json:"width" bson:"width"`
	Height    int       `json:"height" bson:"height"`
	URL       string    `json:"url" bson:"url"`
	IsSolid   bool      `json:"is_solid" bson:"is_solid"`
	CreatedBy string    `json:"created_by" bson:"created_by"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
