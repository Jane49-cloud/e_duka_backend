package models

import (
	"eleliafrika.com/backend/database"
	"gorm.io/gorm"
)

type SubCategory struct {
	gorm.Model
	SubCategoryID            string `gorm:"primary_key;column:subcategory_id;not null;" json:"subcategoryid"`
	SubCategoryName          string `gorm:"column:subcategory_name;not null;" json:"subcategoryname"`
	SubCategoryImage         string `gorm:"" json:"subcategoryimage"`
	SubCategoryTotalProducts int32  `gorm:"default:0" json:"subcatproducts"`
	IsDeleted                bool   `gorm:"column:is_deleted;default:false" json:"isdeleted"`
	ParentCategory           string `gorm:"column:parent_category" json:"parentcategory"`
}

func (subcategory *SubCategory) Save() (*SubCategory, error) {
	err := database.Database.Create(&subcategory).Error
	if err != nil {
		return &SubCategory{}, err
	}
	return subcategory, nil
}
