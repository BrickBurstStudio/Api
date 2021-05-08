package model

import (
	guuid "github.com/google/uuid"
)

type Hub struct {
	ID 	  	  guuid.UUID `gorm:"primaryKey" json:"ID"`
	Name	  string	 `json:"name"`
	Value     string 	 `json:"value"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
}
