package models

import uuid "github.com/satori/go.uuid"

type UserRegisterModel struct {
	ObjectId        uuid.UUID `json:"objectId"`
	Username        string    `bson:"username" json:"username"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirmPassword"`
}
