package model

type File struct {
	Url	  		string	 	`gorm:"primaryKey" json:"url"`
	UpdatedAt 	int64     	`gorm:"autoUpdateTime" json:"-"`
}
