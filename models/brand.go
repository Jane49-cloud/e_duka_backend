package models

import (
	"eleliafrika.com/backend/database"
	"gorm.io/gorm"
)

type Brand struct {
	gorm.Model
	BrandID          string `gorm:"primary_key;column:brand_id;not null;unique" json:"brandid"`
	BrandName        string `gorm:"column:brand_name;unique;not null" json:"brandname"`
	Imageurl         string `gorm:"column:image_url" json:"imageurl"`
	Isdeleted        bool   `gorm:"column:is_deleted;default:false" json:"isdeleted"`
	TotalProducts    int    `gorm:"default:0;column:total_products;" json:"totalproducts"`
	TotalEngagements int    `gorm:"default:0;column:total_engagements" json:"totalengagements"`
}

func (brand *Brand) Save() (*Brand, error) {
	err := database.Database.Create(&brand).Error
	if err != nil {
		return &Brand{}, nil
	}
	return brand, nil
}
