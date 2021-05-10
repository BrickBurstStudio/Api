package model

import (
	"time"

	guuid "github.com/google/uuid"
)

type Key struct {
	ID 	  	  	guuid.UUID 	`gorm:"primaryKey" json:"id"`
	Expires   	time.Time  	`json:"expires"`
	IP			string		`json:"ip"`
	CreatedAt 	int64      	`gorm:"autoCreateTime" json:"-" `
	UpdatedAt 	int64      	`gorm:"autoUpdateTime" json:"updated"`
	Check1 		bool		`json:"c1"`
	Check2 		bool		`json:"c2"`
	Check3 		bool		`json:"c3"`
	Check4 		bool		`json:"c4"`
	Check5 		bool		`json:"c5"`
}
