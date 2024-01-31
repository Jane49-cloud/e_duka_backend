package product

import (
	"eleliafrika.com/backend/database"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ProductID          string `gorm:"column:product_id;not null;primary key;unique;" json:"producttid"`
	ProductName        string `gorm:"column:product_name;not null" json:"productname"`
	ProductPrice       string `gorm:"column:product_price;not null" json:"productprice"`
	ProductDescription string `gorm:"column:product_description;" json:"productdescription"`
	UserID             string `gorm:"size:255;not null;" json:"userid"`
	MainImage          string `gorm:"type:text;size:65535;" json:"mainimage"`
	IsSuspended        bool   `gorm:"column:is_suspended;default:false;not null;" json:"issuspended"`
	IsApproved         bool   `gorm:"column:is_approved;default:false;not null;" json:"isapproved"`
	Quantity           int    `gorm:"default:0" json:"quantity"`
	IsActive           bool   `gorm:"column:is_active;default:true" json:"isactive"`
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

type ProductCount struct {
	gorm.Model
	Username          string `gorm:"column:username;not null;" json:"username"`
	UserID            string `gorm:"column:user_id;not null;" json:"userId"`
	TotalProducts     uint   `gorm:"column:total_products;default:0;" json:"totalProducts"`
	ActiveProducts    uint   `gorm:"column:active_products;default:0;" json:"activeProducts"`
	InActiveProducts  uint   `gorm:"column:in_active_products;default:0;" json:"inActiveProducts"`
	DeletedProducts   uint   `gorm:"column:deleted_products;default:0;" json:"deletedProducts"`
	SuspendedProducts uint   `gorm:"column:suspended_reached;default:0;" json:"suspendedProducts"`
	LimitReached      bool   `gorm:"column:limit_reached;default:false;" json:"limitReached"`
}

func (productCount *ProductCount) Save() (*ProductCount, error) {
	err := database.Database.Create(&productCount).Error
	if err != nil {
		return &ProductCount{}, err
	}
	return productCount, nil
}
