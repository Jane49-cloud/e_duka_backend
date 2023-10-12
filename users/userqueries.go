package users

import (
	"errors"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
)

// query user using their email
func FindUserByEmail(email string) (models.User, error) {
	if len(email) < 10 {
		return models.User{}, errors.New("user email provided is null")
	}
	var user models.User
	err := database.Database.Where("email=?", email).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// function to query user with id
func FindSellerById(id string) (models.User, error) {
	var user models.User
	err := database.Database.Where("user_id=?", id).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
