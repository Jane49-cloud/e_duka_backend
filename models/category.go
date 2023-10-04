package models

import (
	"eleliafrika.com/backend/database"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	CategoryID    string `gorm:"primary_key;column:category_id;not null;" json:"categoryid"`
	CategoryName  string `gorm:"column:category_name;not null;" json:"categoryname"`
	CategoryImage string `gorm:"" json:"categoryimage"`
	TotalProducts int32  `gorm:"default:0" json:"totalproducts"`
	TotalRevenue  int32  `gorm:"default:0" json:"totalrevenue"`
}

func (category *Category) Save() (*Category, error) {
	err := database.Database.Create(&category).Error
	if err != nil {
		return &Category{}, err
	}
	return category, nil
}
