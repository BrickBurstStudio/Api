package model

type File struct {
	Url	  		string	 	`json:"url"`
	UpdatedAt 	int64     	`gorm:"autoUpdateTime" json:"-"`
}
