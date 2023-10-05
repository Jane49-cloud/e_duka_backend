package users

import (
	"net/http"

	"time"

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
		response := models.Reply{
			Message: "Error validating user",
			Error:   err.Error(),
			Success: false,
			Data:    input,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
	if !success {
		response := models.Reply{
			Message: "Error validating user",
			Success: false,
			Data:    input,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	// get current data to save user with
	randomuuid := uuid.New()
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	user := models.User{
		UserID:          randomuuid.String(),
		Firstname:       input.Firstname,
		Middlename:      input.Middlename,
		Lastname:        input.Lastname,
		UserImage:       input.UserImage,
		UserLocation:    input.UserLocation,
		Email:           input.Email,
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

		savedUser, err := user.Save()
		if err != nil {
			response := models.Reply{
				Message: "Error creating user",
				Error:   err.Error(),
				Success: false,
				Data:    user,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}
		// generate token directly on succesfuly register
		token, err := GenerateJWT(user)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "could not generate token for the user",
			})
		}

		context.JSON(http.StatusCreated, gin.H{
			"user":    savedUser,
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
		response := models.Reply{
			Message: "Error validating user input",
			Error:   err.Error(),
			Success: false,
			Data:    input,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if !success {
		response := models.Reply{
			Message: "Error validating user",
			Error:   err.Error(),
			Success: false,
			Data:    input,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {

		// check if user exists

		user, err := FindUserByEmail(input.Email)

		if err != nil {
			response := models.Reply{
				Message: "Error fetching user",
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
				Message: "incorrect password",
				Success: false,
				Error:   err.Error()}
			context.JSON(http.StatusBadRequest, response)
			return
		}

		// generate jwt if error does not exists
		token, err := GenerateJWT(user)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Error occured while generating the token",
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
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusOK, response)
			return
		} else {
			userData := models.User{
				Firstname:    user.Firstname,
				Middlename:   user.Middlename,
				Lastname:     user.Lastname,
				Email:        user.Email,
				UserImage:    user.UserImage,
				UserLocation: user.UserLocation,
				UserID:       user.UserID,
			}
			response := models.Reply{
				Message: "Succesfully fetched the user",
				Success: true,
				Data:    userData,
			}
			context.JSON(http.StatusOK, response)
			return
		}
	}
}
