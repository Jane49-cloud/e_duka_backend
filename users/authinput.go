package users

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
