package users

import "reflect"

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

func TestNullOrEmpty(register RegisterInput) bool {
	registerType := reflect.TypeOf(register)
	registerValue := reflect.ValueOf(register)

	for i := 0; i < registerType.NumField(); i++ {
		fieldValue := registerValue.Field(i)

		// Define your null criteria (e.g., empty string for strings, zero for numeric types)
		nullCriteria := reflect.Zero(fieldValue.Type())

		// Compare the field value with the null criteria
		if reflect.DeepEqual(fieldValue.Interface(), nullCriteria.Interface()) {
			return true // Field has a "null" value
		}
	}

	return false
}
