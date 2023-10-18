package admin

import (
	"errors"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/images"
	"eleliafrika.com/backend/users"
)

func ValidateRegisterInput(admin *SystemAdmin) (bool, error) {
	// details := []string{admin.AdminName, admin.Email, admin.Cell, admin.Password, admin.AdminImage, admin.Role}
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
