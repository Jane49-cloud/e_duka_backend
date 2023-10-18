package admin

import (
	"html"
	"strings"

	"eleliafrika.com/backend/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SystemAdmin struct {
	gorm.Model
	AdminID       string `gorm:"not null;primary_key;unique" json:"userid"`
	AdminName     string `gorm:"size:255;not null;column:admin_name;" json:"adminname"`
	Email         string `gorm:"size:255;not null;unique;column:email" json:"email"`
	Cell          string `gorm:"size:255;not null;unique;column:cell" json:"cell"`
	Password      string `gorm:"size:255;not null;" json:"password"`
	AdminImage    string `gorm:"size:255;column:admin_image;" json:"adminimage"`
	Role          string `gorm:"column:role;not null;default:'basic';" json:"role"`
	DateAdded     string `gorm:"column:date_added;" json:"dateadded"`
	Token         string `gorm:"column:admin_token;" json:"token"`
	LastLoggedIn  string `gorm:"column:last_logged_in;" json:"lastlogin"`
	Notifications int    `gorm:"column:notifications;default:0" json:"notifications"`
	Chats         int    `gorm:"column:chats;default:0;" json:"chats"`
}

type AddAdmin struct {
	AdminName  string `json:"adminname"`
	Email      string `json:"email"`
	Cell       string `json:"cell"`
	Password   string `json:"password"`
	AdminImage string `json:"adminimage"`
	Role       string `json:"role"`
}

type AdminLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// function to create new user
func (systemadmin *SystemAdmin) Save() (*SystemAdmin, error) {

	err := database.Database.Create(&systemadmin).Error
	if err != nil {
		return &SystemAdmin{}, err
	}
	return &SystemAdmin{}, nil

}

// hash the password before saving the user in the database
func (systemadmin *SystemAdmin) BeforeSave(*gorm.DB) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(systemadmin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err

	}

	systemadmin.Password = string(hashPassword)
	systemadmin.AdminName = html.EscapeString(strings.TrimSpace(systemadmin.AdminName))
	systemadmin.Email = html.EscapeString(strings.TrimSpace(systemadmin.Email))
	return nil
}
