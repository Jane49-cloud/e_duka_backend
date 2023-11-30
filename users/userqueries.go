package users

import (
	"errors"

	"eleliafrika.com/backend/database"
)

// query user using their email
func FindUserByEmail(email string) (User, error) {
	if len(email) < 10 {
		return User{}, errors.New("user email provided is null")
	}
	var user User
	err := database.Database.Where("email=?", email).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// query user using their email
func FindUserByPhone(phone string) (User, error) {
	if len(phone) < 10 {
		return User{}, errors.New("phone email provided is null")
	}
	var user User
	err := database.Database.Where("phone=?", phone).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// function to query user with id
func FindUserById(id string) (User, error) {
	var user User
	err := database.Database.Where("user_id=?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}
