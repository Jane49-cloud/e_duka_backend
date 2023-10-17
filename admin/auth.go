package admin

import (
	"strings"

	"eleliafrika.com/backend/database"
	"golang.org/x/crypto/bcrypt"
)

func (admin *SystemAdmin) CreateAdmin() (*SystemAdmin, error) {

	err := database.Database.Create(&admin).Error
	if err != nil {
		return &SystemAdmin{}, err
	}
	return &SystemAdmin{}, nil
}

func (admin *SystemAdmin) BeforeCreate() error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin.Password = string(encryptedPassword)
	admin.AdminName = strings.TrimSpace(admin.AdminName)
	admin.Cell = strings.ReplaceAll(admin.Cell, " ", "")
	admin.Email = strings.TrimSpace(admin.Email)
	admin.Role = strings.TrimSpace(admin.Role)

	return nil
}

func ValidateHashPassword(hashedPassword string, userPassword string) (bool, error) {
	if userPassword == "" {
		return false, bcrypt.ErrMismatchedHashAndPassword

	}

	if hashedPassword == "" {
		return false, bcrypt.ErrHashTooShort

	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))
	if err != nil {
		return false, err
	}
	return true, nil
}
