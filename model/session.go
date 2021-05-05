package model

import (
	"time"

	guuid "github.com/google/uuid"
)

type Session struct {
	SessionID guuid.UUID `gorm:"primaryKey" json:"sessionID"`
	Expires   time.Time  `json:"-"`
	UserRefer guuid.UUID `json:"-"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
}
