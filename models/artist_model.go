package models

type Artist struct {
	Base
	Name string `json:"name" gorm:"type:varchar(150)"`
}
