package product

import "gorm.io/gorm"

type AddProductInput struct {
	gorm.Model
	ProductName        string   `gorm:"column:product_name;unique;not null" json:"productname"`
	ProductPrice       string   `gorm:"column:product_price;not null" json:"productprice"`
	ProductDescription string   `gorm:"column:product_description;" json:"productdescription"`
	MainImage          string   `gorm:"not null;" json:"mainimage"`
	ProductImages      []string `gorm:"type:text[]" json:"productimages"`
	Quantity           int      `gorm:"default:0" json:"quantity"`
	ProductType        string   `gorm:"column:product_type;" json:"producttype"`
	Brand              string   `gorm:"column:brand" json:"brand"`
	Category           string   `gorm:"category" json:"category"`
	SubCategory        string   `gorm:"column:subcategory" json:"subcategory"`
}
