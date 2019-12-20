package models

import uuid "github.com/satori/go.uuid"

type UpdateSettingGroupModel struct {
	Type        string                     `json:"type"`
	CreatedDate int64                      `json:"created_date"`
	OwnerUserId uuid.UUID                  `json:"ownerUserId"`
	List        []GetSettingGroupItemModel `json:"list"`
}
