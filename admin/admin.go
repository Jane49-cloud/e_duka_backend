package admin

import "gorm.io/gorm"

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
