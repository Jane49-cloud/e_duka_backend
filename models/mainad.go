package models

import "eleliafrika.com/backend/database"

type MainAd struct {
	Advertid    string `gorm:"column:ad_id;primary key;unique;not null" json:"advertid"`
	AdBy        string `gorm:"column:ad_by;not null;type:text" json:"adby"`
	AdName      string `gorm:"column:ad_name;not null;type:text" json:"adname"`
	AdImage     string `gorm:"column:ad_image;not null;type:text" json:"adimage"`
	IsActive    bool   `gorm:"column:is_active;default:true" json:"isactive"`
	IsDeleted   bool   `gorm:"column:is_deleted;default:false;" json:"isdeleted"`
	DateCreated string `gorm:"column:created_on;not null" json:"datecreated"`
	EndingDate  string `gorm:"column:ending_date;not null;" json:"endingdate"`
	AdCategory  string `gorm:"column:ad_category;type:text;not null" json:"adcategory"`
}

func (mainad *MainAd) Save() (*MainAd, error) {
	err := database.Database.Create(&mainad).Error
	if err != nil {
		return &MainAd{}, err
	}

	return mainad, nil
}
