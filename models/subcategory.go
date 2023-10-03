package models

import "gorm.io/gorm"

type SubCategory struct {
	gorm.Model
	SubCategoryID            string   `gorm:"primary_key;column:sub_category_id;not null;" json:"subcategoryid"`
	SubCategoryName          string   `gorm:"column:sub_category_name;not null;" json:"subcategoryname"`
	SubCategoryImage         string   `gorm:"" json:"subcategorycategoryimage"`
	SubCategoryTotalProducts int32    `gorm:"default:0" json:"subcategorytotalproducts"`
	Products                 []string `gorm:"column:products;type:jsonb" json:"products"`
	CategoryID               string   `gorm:"column:category_id" json:"categoryid"`
	Category                 string   `gorm:"column:category" json:"parentcategory"`
}
