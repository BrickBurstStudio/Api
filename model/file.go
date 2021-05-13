package model

import (
	guuid "github.com/google/uuid"
)

type File struct {
	ID			guuid.UUID	`gorm:"primary_key" json:"id"`
	Url	  		string	 	`json:"url"`
	UpdatedAt 	int64     	`gorm:"autoUpdateTime" json:"-"`
}
