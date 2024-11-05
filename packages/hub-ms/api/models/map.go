package models

import "time"

type Map struct {
	ID          string                           `json:"id" bson:"_id"`
	Name        string                           `json:"name" bson:"name"`
	Width       int                              `json:"width" bson:"width"`
	Height      int                              `json:"height" bson:"height"`
	MapElements map[string]map[string]MapElement `json:"map_elements" bson:"map_elements"` // key will be either 'true' (invalidPos elements) or 'false'(validPos elements) and secondary key will be "x_y" depicting the coordinates
	CreatedBy   string                           `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time                        `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time                        `json:"updated_at" bson:"updated_at"`
}

type MapElement struct {
	ID    string `json:"id" bson:"_id"`
	EleId string `json:"ele_id" bson:"ele_id"`
}
