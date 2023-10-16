package product

import (
	"eleliafrika.com/backend/database"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ProductID          string `gorm:"column:product_id;not null;primary key;unique;" json:"producttid"`
	ProductName        string `gorm:"column:product_name;unique;not null" json:"productname"`
	ProductPrice       string `gorm:"column:product_price;not null" json:"productprice"`
	ProductDescription string `gorm:"column:product_description;" json:"productdescription"`
	UserID             string `gorm:"size:255;not null;" json:"userid"`
	MainImage          string `gorm:"not null;" json:"mainimage"`
	IsSuspended        bool   `gorm:"column:is_suspended;default:false;not null;" json:"issuspended"`
	IsApproved         bool   `gorm:"column:is_approved;default:false;not null;" json:"isapproved"`
	Quantity           int    `gorm:"default:0" json:"quantity"`
	IsActive           bool   `gorm:"column:is_active;default:false" json:"isactive"`
	IsDeleted          bool   `gorm:"column:is_deleted;default:false" json:"isdeleted"`
	ActiveUntil        string `gorm:"column:active_until" json:"activeuntil"`
	ProductType        string `gorm:"column:product_type;" json:"producttype"`
	TotalLikes         int    `gorm:"default:0" json:"totallikes"`
	TotalComments      int    `gorm:"default:0" json:"totalcomments"`
	DateAdded          string `gorm:"" json:"dateadded"`
	LastUpdated        string `gorm:"size:255;not null" json:"lastupdated"`
	LatestInteractions string `gorm:"size:255;not null" json:"latestinteractions"`
	TotalInteractions  int    `gorm:"size:255;not null" json:"totalinteractions"`
	TotalBookmarks     int    `gorm:"size:255;not null" json:"totalbookmarks"`
	Brand              string `gorm:"column:brand" json:"brand"`
	Category           string `gorm:"category" json:"category"`
	SubCategory        string `gorm:"column:subcategory" json:"subcategory"`
}

type AddProductInput struct {
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

func (product *Product) Save() (*Product, error) {
	err := database.Database.Create(&product).Error
	if err != nil {
		return &Product{}, err
	}
	return product, nil
}
