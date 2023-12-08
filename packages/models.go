package packages

import (
	"eleliafrika.com/backend/database"
	"gorm.io/gorm"
)

type PackageModel struct {
	gorm.Model
	PackageId   string `gorm:"column:package_id;not null;" json:"package_id"`
	PackageName string `gorm:"column:package_name;not null;" json:"package_name"`
	UsersNumber uint   `gorm:"column:users_number;not null;" json:"users_number"`
	Price       uint   `gorm:"column:price;not null;" json:"price"`
	Duration    int    `gorm:"column:duration;not null;" json:"duration"`
	DateCreated string `gorm:"column:date_created;not null;" json:"date_created"`
	DateUpdated string `gorm:"column:date_updated;not null;" json:"date_updated"`
}

type AddPackage struct {
	PackageName string `gorm:"column:package_name;not null;" json:"package_name"`
	Price       uint   `gorm:"column:price;not null;" json:"price"`
	Duration    int    `gorm:"column:duration;not null;" json:"duration"`
}

func (packagemodel *PackageModel) Save() (*PackageModel, error) {
	err := database.Database.Create(&packagemodel).Error
	if err != nil {
		return &PackageModel{}, err
	}
	return packagemodel, nil
}
