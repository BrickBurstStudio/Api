package model

import (
	guuid "github.com/google/uuid"
)

type File struct {
	Url	  		string	 	`json:"url"`
	ID			guuid.UUID	`gorm:"primary_key" json:"id"`
	UpdatedAt 	int64     	`gorm:"autoUpdateTime" json:"-"`
}
