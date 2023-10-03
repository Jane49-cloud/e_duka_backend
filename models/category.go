package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	CategoryID    string   `gorm:"primary_key;column:category_id;not null;" json:"categoryid"`
	CategoryName  string   `gorm:"column:category_name;not null;" json:"categoryname"`
	CategoryImage string   `gorm:"" json:"categoryimage"`
	TotalProducts int32    `gorm:"default:0" json:"totalproducts"`
	Products      []string `gorm:"column:products;type:jsonb" json:"products"`
	Brands        []string `gorm:"column:brands;type:jsonb" json:"brands"`
	SubCategories []string `gorm:"column:sub_categories;type:jsonb" json:"subcategories"`
}
