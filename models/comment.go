package models

import (
	"eleliafrika.com/backend/database"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	CommentID     string `gorm:"not null;size:255;column:comment_id;primary_key" json:"commentid"`
	ProductID     string `gorm:"not null;size:255;column:product_id;primary_key" json:"productid"`
	UserID        string `gorm:"not null;size:255" json:"userid"`
	Comment       string `gorm:"not null;size:255;type:text;" json:"comment"`
	Isdeleted     bool   `gorm:"default:false;type:bool" json:"isdeleted"`
	DateCommented string `json:"datecommented"`
}

type Commentinput struct {
	gorm.Model
	CommentID string `json:"commentid"`
	ProductID string `json:"productid"`
	UserID    string `json:"userid"`
	Comment   string `json:"comment"`
	Isdeleted bool   `json:"isdeleted"`
}

func (comment *Comment) Save() (*Comment, error) {
	err := database.Database.Create(&comment).Error
	if err != nil {
		return &Comment{}, err
	}
	return comment, nil
}
