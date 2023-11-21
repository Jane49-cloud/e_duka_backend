package users

import (
	"errors"
	"net/http"
	"strings"

	"time"

	"eleliafrika.com/backend/images"
	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(context *gin.Context) {
	var input RegisterInput

	if err := context.ShouldBindJSON(&input); err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error binding data",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
	success, err := ValidateRegisterInput(&input)
	if err != nil {

		response := models.Reply{
			Error:   err.Error(),
			Message: "error validating data",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
	if !success {

		response := models.Reply{
			Error:   errors.New("returned false").Error(),
			Message: "error validating data",
			Success: false,
		}

		context.JSON(http.StatusBadRequest, response)
		return
	}

	// get current data to save user with
	randomuuid := uuid.New()
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	imageUrl, err := images.UploadHandler(input.Firstname+input.Lastname, input.UserImage, context)
	if err != nil {
		response := models.Reply{
			Message: "main image not saved",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	user := User{
		UserID:          randomuuid.String(),
		Firstname:       input.Firstname,
		Middlename:      input.Middlename,
		Lastname:        input.Lastname,
		Location:        input.UserLocation,
		Email:           input.Email,
		UserImage:       imageUrl,
		Phone:           input.Phone,
		Password:        input.Password,
		DateJoined:      formattedTime,
		LastLoggedIn:    formattedTime,
		LastInteraction: formattedTime,
	}

	userExists, err := FindUserByEmail(user.Email)

	if err != nil {

		response := models.Reply{
			Error:   err.Error(),
			Message: "error fetching user data",
			Success: false,
		}

		context.JSON(http.StatusBadRequest, response)
		return
	} else if userExists.Email != "" {

		response := models.Reply{
			Error:   errors.New("user does exist").Error(),
			Message: "email has already been used",
			Success: false,
		}

		context.JSON(http.StatusBadRequest, response)
		return
	} else {

		_, err := user.Save()
		if err != nil {

			response := models.Reply{
				Error:   err.Error(),
				Message: "error creating user",
				Success: false,
			}

			context.JSON(http.StatusBadRequest, response)
			return
		}
		// generate token directly on succesfuly register
		token, err := GenerateJWT(user)

		if err != nil {
			response := models.Reply{
				Message: "could not generate token for the user",
				Success: false,
				Data:    token,
				Error:   err.Error(),
			}
			context.JSON(http.StatusBadRequest, response)
		}

		response := models.Reply{
			Message: "User has been created succesfully",
			Success: true,
			Data:    token,
		}
		context.JSON(http.StatusCreated, response)
	}

}

// code for login in
func Login(context *gin.Context) {
	var input LoginInput

	if err := context.ShouldBindJSON(&input); err != nil {
		response := models.Reply{
			Message: "error binding the user input",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
	// check validity of user input
	success, err := ValidateLoginInput(&input)
	if err != nil {
		response := models.Reply{
			Message: "error validating user input",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if !success {

		response := models.Reply{
			Message: "could not validate date",
			Error:   errors.New("validation returned false").Error(),
			Success: false,
		}

		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		// check if user exists
		user, err := FindUserByEmail(input.Email)

		if err != nil {
			response := models.Reply{
				Message: "error fetching user",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}
		// validate the password password passed with the harsh on db
		err = user.ValidatePassword(input.Password)
		if err != nil {
			response := models.Reply{
				Message: "incorrect details",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}
		// generate jwt if error does not exists
		token, err := GenerateJWT(user)

		if err != nil {
			response := models.Reply{
				Message: "error occured on authentication",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return

		}

		response := models.Reply{
			Message: "Authentication successful",
			Data:    token,
			Success: true,
		}
		context.JSON(http.StatusOK, response)
	}
}

func GetSingleUser(context *gin.Context) {
	user, err := CurrentUser(context)
	if err != nil {

		response := models.Reply{
			Message: "error fetching current user",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		if user.Firstname == "" {

			response := models.Reply{
				Message: "user does not exist",
				Error:   errors.New("error user does not exist").Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			userData := User{
				Firstname:  user.Firstname,
				Middlename: user.Middlename,
				Lastname:   user.Lastname,
				Email:      user.Email,
				UserImage:  user.UserImage,
				Location:   user.Location,
				UserID:     user.UserID,
				IsApproved: user.IsApproved,
				Phone:      user.Phone,
			}

			response := models.Reply{
				Message: "Succesfully fetched the user",
				Data:    userData,
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		}
	}
}
func FetchSingleUser(context *gin.Context) {
	id := context.Query("id")
	user, err := FindUserById(strings.ReplaceAll(id, "'", ""))
	if err != nil {
		response := models.Reply{
			Message: "error fetching user",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if user.Firstname == "" {
		response := models.Reply{
			Message: "user does not exist",
			Error:   errors.New("could not find user").Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		if err != nil {
			response := models.Reply{
				Message: "could not fetch the user image",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}
		response := models.Reply{
			Message: "user fetched succesfully",
			Data:    user,
			Success: true,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
}
func UpdateUser(context *gin.Context) {

	var userUpdateData User
	if err := context.ShouldBindJSON(&userUpdateData); err != nil {
		response := models.Reply{
			Message: "could not bind the user data to the request needs",
			Error:   err.Error(),
			Success: false,
			Data:    userUpdateData,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
	thisUser, err := CurrentUser(context)
	if err != nil {
		response := models.Reply{
			Message: "could not fetch current user",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if thisUser.Firstname == "" {
		response := models.Reply{
			Message: "user not found",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		userid := context.Query("userid")

		query := strings.ReplaceAll(userid, "'", "")
		newUser := User{
			Firstname:  userUpdateData.Firstname,
			Middlename: userUpdateData.Middlename,
			Lastname:   userUpdateData.Lastname,
			UserImage:  userUpdateData.UserImage,
			Location:   userUpdateData.Location,
			Email:      userUpdateData.Email,
			Phone:      userUpdateData.Phone,
		}
		updateUser, err := UpdateUserUtil(query, newUser)
		if err != nil {
			response := models.Reply{
				Message: "could not update user",
				Success: false,
				Error:   err.Error(),
			}
			context.JSON(http.StatusOK, response)
			return
		} else {
			response := models.Reply{
				Message: "user updated successfully",
				Success: true,
				Data:    updateUser,
			}
			context.JSON(http.StatusOK, response)
			return
		}
	}
}
func FetchSellers(context *gin.Context) {
	users, err := FetchAllSellersUtil()

	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error fetching all users",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		response := models.Reply{
			Data:    users,
			Message: "succesfully fetched all users",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	}
}
