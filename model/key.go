package model

import (
	"time"

	guuid "github.com/google/uuid"
)

type Key struct {
	ID 	  	  guuid.UUID `gorm:"primaryKey" json:"ID"`
	Expires   time.Time  `json:"expires"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
}
