package models

import "time"

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Username  string    `json:"username" bson:"username"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	AvatarId  string    `json:"avatarId" bson:"avatarId"`
	Role      string    `json:"role" bson:"role"`
	Rooms     []string  `json:"rooms" bson:"rooms"`
	Maps      []string  `json:"maps" bson:"maps"`
	Elements  []string  `json:"elements" bson:"elements"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type Avatar struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}
