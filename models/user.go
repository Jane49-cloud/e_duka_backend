package models

import (
	"html"
	"strings"

	"eleliafrika.com/backend/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID                  string   `gorm:"not null;primary_key;unique" json:"userid"`
	Firstname               string   `gorm:"size:255;not null" json:"username"`
	Middlename              string   `gorm:"size:255;not null" json:"middlename"`
	Lastname                string   `gorm:"size:255;not null" json:"lastname"`
	UserImage               string   `gorm:"size:255;" json:"userimage"`
	UserLocation            string   `gorm:"size:255;not null" json:"location"`
	Email                   string   `gorm:"size:255;not null;unique" json:"email"`
	Phone                   string   `gorm:"size:255;not null;unique" json:"phone"`
	Password                string   `gorm:"size:255;not null;" json:"password"`
	Totalproducts           int      `gorm:"default:0;column:total_products;default:0" json:"totalproducts"`
	Products                []string `gorm:"column:products;type:jsonb;default:null" json:"products"`
	UserType                string   `gorm:"column:user_type;not null;default:'visitor';" json:"usertype"`
	TotalLikes              int      `gorm:"default:0;column:total_likes;" json:"totallikes"`
	TotalViews              int      `gorm:"default:0;column:total_views;" json:"totalviews"`
	TotalEngagements        int      `gorm:"default:0;column:total_engagements;" json:"totalengagements"`
	Followers               int      `gorm:"default:0;column:followers;" json:"followers"`
	Follows                 int      `gorm:"default:0;column:follows;" json:"follows"`
	DateJoined              string   `gorm:"column:date_joined;" json:"datejoined"`
	LastLoggedIn            string   `gorm:"column:last_logged_in;" json:"lastlogin"`
	LastInteraction         string   `gorm:"column:last_interaction;" json:"lastinteraction"`
	LastproductAdd          string   `gorm:"column:last_product_add;" json:"lastproductadd"`
	TotalActiveproducts     int      `gorm:"default:0;column:total_active_products;" json:"totalactiveproducts"`
	DailyAverageEngagements string   `gorm:"default:0;column:daily_average_engagements;" json:"dailyaverageengagements"`
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
