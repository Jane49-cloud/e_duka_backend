package users

import (
	"errors"
	"net/http"
	"strings"

	"time"

	globalutils "eleliafrika.com/backend/global_utils"
	"eleliafrika.com/backend/images"
	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(context *gin.Context) {
	var input RegisterInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	success, err := ValidateRegisterInput(&input)
	if err != nil {

		context.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "error validating user",
			"token":   nil,
		})
		return
	}
	if !success {

		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "error validating user",
			"token":   nil,
		})
		return
	}

	// get current data to save user with
	randomuuid := uuid.New()
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	userImagepath, err := UploadUserImage(input.UserImage, input.Firstname+"-"+input.Lastname)

	if err != nil {
		globalutils.HandleError("error uploading user image", err, context)
		return
	}

	user := User{
		UserID:          randomuuid.String(),
		Firstname:       input.Firstname,
		Middlename:      input.Middlename,
		Lastname:        input.Lastname,
		Location:        input.UserLocation,
		Email:           input.Email,
		UserImage:       userImagepath,
		Phone:           input.Phone,
		Password:        input.Password,
		DateJoined:      formattedTime,
		LastLoggedIn:    formattedTime,
		LastInteraction: formattedTime,
	}

	userExists, err := FindUserByEmail(user.Email)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	} else if userExists.Email != "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "email has already been used",
		})
	} else {

		_, err := user.Save()
		if err != nil {

			context.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "error creating user",
				"token":   nil,
			})
			return
		}
		// generate token directly on succesfuly register
		token, err := GenerateJWT(user)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "could not generate token for the user",
				"token":   nil,
			})
		}

		context.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "User has been created succesfully",
			"token":   token,
		})
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
		context.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "error validating user input",
			"token":   nil,
		})
		return
	} else if !success {

		context.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "error validating user",
			"token":   nil,
		})
		return
	} else {

		// check if user exists

		user, err := FindUserByEmail(input.Email)

		if err != nil {

			context.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "error fetching user",
				"token":   nil,
			})
			return
		}

		// validate the password password passed with the harsh on db

		err = user.ValidatePassword(input.Password)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "incorrect details",
				"token":   nil,
			})
			return
		}

		// generate jwt if error does not exists
		token, err := GenerateJWT(user)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Error occured while generating the token",
				"token":   nil,
			})
			return

		}

		context.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Authentication successful",
			"token":   token,
		})
	}
}

func GetSingleUser(context *gin.Context) {
	user, err := CurrentUser(context)
	if err != nil {

		globalutils.HandleError("error fetching current user", err, context)
		return
	} else {
		if user.Firstname == "" {
			globalutils.HandleError("user does not exist", err, context)
			return
		} else {

			userImage, err := images.DownloadImageFromBucket(user.UserImage)
			if err != nil {
				globalutils.HandleError("could not fetch the user image", err, context)
			}
			userData := User{
				Firstname:  user.Firstname,
				Middlename: user.Middlename,
				Lastname:   user.Lastname,
				Email:      user.Email,
				UserImage:  userImage,
				Location:   user.Location,
				UserID:     user.UserID,
				IsApproved: user.IsApproved,
			}

			globalutils.HandleSuccess("Succesfully fetched the user", userData, context)
			return
		}
	}
}
func FetchSingleUser(context *gin.Context) {
	id := context.Query("id")
	user, err := FindUserById(strings.ReplaceAll(id, "'", ""))
	if err != nil {
		globalutils.HandleError("error fetching user", err, context)
		return
	} else if user.Firstname == "" {
		globalutils.HandleError("user does not exist", errors.New("user does not exist"), context)
		return
	} else {
		userImage, err := images.DownloadImageFromBucket(user.UserImage)
		if err != nil {
			globalutils.HandleError("could not fetch the user image", err, context)
		}
		user.UserImage = userImage
		globalutils.HandleSuccess("user feteched succesfully", user, context)
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

			globalutils.HandleError("could not update user", err, context)
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
		globalutils.HandleError("error fetching all users", err, context)
		return
	} else {

		globalutils.HandleSuccess("succesfully fetched all users", users, context)
		return
	}
}
