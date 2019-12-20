package models

import uuid "github.com/satori/go.uuid"

type GetSettingGroupItemModel struct {
	ObjectId uuid.UUID `json:"objectId"`
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	IsSystem bool      `json:"isSystem"`
}
