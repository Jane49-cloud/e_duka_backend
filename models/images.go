package models

import (
	"fmt"

	"eleliafrika.com/backend/database"
	"gorm.io/gorm"
)

type ProductImage struct {
	gorm.Model
	ImageID   string `gorm:"primary_key;not null;unique"`
	ProductID string `gorm:"not null" json:"productid"`
	ImageUrl  string `gorm:"type:text;size:65535;not null;" json:"imageurl"`
}

func (image *ProductImage) Save() (*ProductImage, error) {
	err := database.Database.Create(&image).Error
	fmt.Printf("data %v", image)
	if err != nil {
		return &ProductImage{}, err
	}
	return image, nil
}
