package admin

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/images"
	"eleliafrika.com/backend/product"
	"eleliafrika.com/backend/users"
)

func ValidateRegisterInput(admin *AddAdmin) (bool, error) {
	details := []string{admin.AdminName, admin.Email, admin.Cell, admin.Password, admin.Role}

	charPattern := "[!@#$%^&*()_+\\-=\\[\\]{};':\"\\\\|,.<>?]"
	emailPattern := "[!#$%^&*()+\\=\\[\\]{};':\"\\\\|,<>?]"
	numPattern := "[0-9]"
	capPattern := "[A-Z]"
	for _, value := range details {
		if value == admin.AdminName {
			value = strings.TrimSpace(admin.AdminName)
			if len(value) < 5 {
				return false, errors.New("admin name is too short")
			} else if len(value) > 20 {
				return false, errors.New("admin name is too long")
			} else if regexp.MustCompile(charPattern).MatchString(admin.AdminName) {
				return false, errors.New("special characters not allowed in admin name")
			}
		} else if value == admin.Email {
			admin.Email = strings.ToLower(admin.Email)
			if len(value) < 8 {
				return false, errors.New("invalid email format!email should contain @, . and should be longer than 8 characters")
			} else if len(value) > 20 {
				return false, errors.New("email too long")
			} else if !strings.Contains(admin.Email, "@") {
				return false, errors.New("invalid email format!email should contain @, . and should be longer than 8 characters")
			} else if !strings.Contains(admin.Email, ".") {
				return false, errors.New("invalid email format!email should contain @, . and should be longer than 8 characters")
			} else if regexp.MustCompile(emailPattern).MatchString(admin.Email) {
				return false, errors.New("invalid characters in email")
			}
		} else if value == admin.Cell {
			if len(value) < 10 {
				return false, errors.New("phone number should be atleast 10 characters long")
			} else if len(value) > 15 {
				return false, errors.New("phone number is too long")
			} else {
				for _, char := range admin.Cell {
					if !unicode.IsNumber(char) {
						return false, errors.New("phone number can only contain numbers")
					}
				}
			}
		} else if value == admin.Password {
			if len(value) < 8 {
				return false, errors.New("password is too short")
			} else if len(value) > 20 {
				return false, errors.New("password is too long")
			} else if !regexp.MustCompile(charPattern).MatchString(admin.Password) {
				return false, errors.New("password must contain atleast one special character")
			} else if !regexp.MustCompile(numPattern).MatchString(admin.Password) {
				return false, errors.New("password must contain atleast numerical digit")
			} else if !regexp.MustCompile(capPattern).MatchString(admin.Password) {
				return false, errors.New("password must contain a capital letter")
			}
		} else if value == admin.Role {
			admin.Role = strings.ToLower(admin.Role)
			if len(value) < 3 {
				return false, errors.New("role string is too short")
			} else if regexp.MustCompile(numPattern).MatchString(admin.Role) {
				return false, errors.New("admin role must not contain a numerical digit")
			}
		}
	}
	return false, nil
}

func FetchAllUsersUtil() ([]users.User, error) {
	var AllUsers []users.User

	err := database.Database.Find(&AllUsers).Error

	if err != nil {
		return []users.User{}, err
	}

	if len(AllUsers) > 0 {
		for _, user := range AllUsers {
			userImage, err := images.DownloadImageFromBucket(user.UserImage)
			if err != nil {
				return []users.User{}, err
			} else if userImage == "" {
				return []users.User{}, errors.New("image not downloaded")
			}
			user.UserImage = userImage
		}
	}

	return AllUsers, nil
}
func ApproveAd(query string) (bool, error) {
	var updatedProduct product.Product
	result := database.Database.Model(&updatedProduct).Where(query).Update("is_approved", true)
	if result.RowsAffected == 0 {
		return false, errors.New("could not deactivate the current product")
	}
	return true, nil
}
