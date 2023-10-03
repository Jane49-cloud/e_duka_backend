package models

import "gorm.io/gorm"

type Brand struct {
	gorm.Model
	BrandID          string   `gorm:"primary_key;column:brand_id;not null;unique" json:"brandid"`
	Imageurl         string   `gorm:"" json:"imageurl"`
	TotalProducts    int      `gorm:"column:total_brand_products;" json:"totalproducts"`
	TotalEngagements int      `gorm:"" json:"totalengagements"`
	Products         []string `gorm:"column:brand_products;type:jsonb" json:"products"`
}
