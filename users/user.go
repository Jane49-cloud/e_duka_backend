package users

import (
	"html"
	"strings"

	"eleliafrika.com/backend/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID          string `gorm:"not null;primary_key;unique" json:"userid"`
	Firstname       string `gorm:"size:255;not null" json:"firstname"`
	Middlename      string `gorm:"size:255;not null" json:"middlename"`
	Lastname        string `gorm:"size:255;not null" json:"lastname"`
	Email           string `gorm:"size:255;not null;unique" json:"email"`
	Phone           string `gorm:"size:255;not null;unique" json:"phone"`
	Password        string `gorm:"size:255;not null;" json:"password"`
	UserImage       string `gorm:"size:255;" json:"userimage"`
	Location        string `gorm:"column:location;size:255;not null" json:"location"`
	NoOfProducts    int    `gorm:"default:0;column:total_products;default:0" json:"noofproducts"`
	PackageType     string `gorm:"column:package_type;not null;default:'basic';" json:"packagetype"`
	ActiveAds       int    `gorm:"column:active_ads;default:0;" json:"activeads"`
	InActiveAds     int    `gorm:"column:in_active_ads;default:0;" json:"inactiveads"`
	DeletedAds      int    `gorm:"column:deleted_ads;default:0;" json:"deletedads"`
	UserType        string `gorm:"column:user_type;not null;default:'visitor';" json:"usertype"`
	IsApproved      bool   `gorm:"column:is_approved;type:bool;default:false;" json:"isapproved"`
	TotalLikes      int    `gorm:"default:0;column:total_likes;" json:"totallikes"`
	TotalViews      int    `gorm:"default:0;column:total_views;" json:"totalviews"`
	DateJoined      string `gorm:"column:date_joined;" json:"datejoined"`
	LastLoggedIn    string `gorm:"column:last_logged_in;" json:"lastlogin"`
	LastInteraction string `gorm:"column:last_interaction;" json:"lastinteraction"`
	Notifications   int    `gorm:"column:notifications;default:0" json:"notifications"`
	Chats           int    `gorm:"column:chats;default:0;" json:"chats"`
	Inquiries       int    `gorm:"column:inquiries;default:0" json:"inquiries"`
}

type RegisterInput struct {
	Firstname    string `gorm:"size:255;not null" json:"firstname"`
	Middlename   string `gorm:"size:255;not null" json:"middlename"`
	Lastname     string `gorm:"size:255;not null" json:"lastname"`
	UserImage    string `gorm:"size:255;" json:"userimage"`
	UserLocation string `gorm:"size:255;not null" json:"location"`
	Email        string `gorm:"size:255;not null;unique" json:"email"`
	Phone        string `gorm:"size:255;not null;unique" json:"phone"`
	Password     string `gorm:"size:255;not null;" json:"password"`
}

type LoginInput struct {
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

// function to create new user
func (user *User) Save() (*User, error) {

	err := database.Database.Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return &User{}, nil

}

// hash the password before saving the user in the database
func (user *User) BeforeSave(*gorm.DB) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err

	}
	user.Password = string(hashPassword)
	user.Firstname = html.EscapeString(strings.TrimSpace(user.Firstname))
	user.Middlename = html.EscapeString(strings.TrimSpace(user.Middlename))
	user.Lastname = html.EscapeString(strings.TrimSpace(user.Lastname))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	return nil
}

// validate the hashed password
func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
