package models

type Artist struct {
	ID   uint   `json:"-" gorm:"primarykey"`
	Name string `json:"name" gorm:"type:varchar(150)"`
	Base
}