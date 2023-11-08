package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	globalutils "eleliafrika.com/backend/global_utils"
	"eleliafrika.com/backend/models"
	"eleliafrika.com/backend/product"
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(context *gin.Context) {
	var input AddAdmin

	if err := context.ShouldBindJSON(&input); err != nil {
		response := models.Reply{
			Error:   errors.New("could not bind data").Error(),
			Message: "could not bind json data from user",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
	success, err := ValidateRegisterInput(&input)
	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error validating user details",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
	if !success {
		response := models.Reply{
			Error:   errors.New("could validate data").Error(),
			Message: "returned false in validating user data",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	// get current data to save user with
	randomuuid := uuid.New()
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	// userImagepath, err := users.UploadUserImage(input.AdminImage, input.AdminName)

	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error uploading user image",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	admin := SystemAdmin{
		AdminID:      randomuuid.String(),
		AdminName:    input.AdminName,
		Email:        input.Email,
		AdminImage:   input.AdminImage,
		Cell:         input.Cell,
		Password:     input.Password,
		DateAdded:    formattedTime,
		LastLoggedIn: formattedTime,
	}

	userExists, err := FindAdminByEmail(admin.Email)

	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error validating user",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if userExists.Email != "" {
		response := models.Reply{
			Error:   errors.New("user email has already been used").Error(),
			Message: "email already used",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {

		_, err := admin.Save()
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
		token, err := GenerateJWT(admin)

		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "error generatin token for user",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}
		response := models.Reply{
			Data:    token,
			Message: "admin added",
			Success: true,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
}

func Login(context *gin.Context) {
	var input AdminLogin

	if err := context.ShouldBindJSON(&input); err != nil {
		response := models.Reply{
			Error:   errors.New("could not bind data").Error(),
			Message: "errror binding data json",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
	// check validity of user input
	success, err := ValidateLoginInput(&input)
	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error validating user input",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if !success {
		response := models.Reply{
			Error:   errors.New("error validating user").Error(),
			Message: "error validating user",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {

		// check if user exists

		admin, err := FindAdminByEmail(input.Email)

		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "error fetching user",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}

		// validate the password password passed with the harsh on db

		err = admin.ValidatePassword(input.Password)
		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "incorect password",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}

		// generate jwt if error does not exists
		token, err := GenerateJWT(admin)

		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "error while generating token",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			response := models.Reply{
				Data:    token,
				Message: "login succesfull",
				Success: true,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}
	}
}

func FetchSellers(context *gin.Context) {

	currentAdmin, err := CurrentUser(context)

	if err != nil {
		globalutils.UnAuthenticated(context)
		return
	} else if currentAdmin.AdminName == "" {
		globalutils.UnAuthorized(context)
		return
	} else {
		users, err := FetchAllUsersUtil()
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
}

func ApproveUser(context *gin.Context) {
	id := context.Query("id")
	id = strings.ReplaceAll(id, "'", "")
	// query := "user_id=" + strings.ReplaceAll(id, "'", "")

	currentAdmin, err := CurrentUser(context)

	if err != nil {
		globalutils.UnAuthenticated(context)
		return
	} else if currentAdmin.AdminName == "" {
		globalutils.UnAuthorized(context)
		return
	} else {

		userExists, err := users.FindUserById(id)
		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "error finding user",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		} else if userExists.Firstname == "" {
			response := models.Reply{
				Error:   errors.New("user cannot be found").Error(),
				Message: "user does not exist",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		} else if userExists.IsApproved {
			response := models.Reply{
				Data:    users.User{},
				Message: "user is already approved",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		} else {

			fmt.Printf("request id\n%v\n", id)
			_, err := users.UpdateUserSpecificField(id, "is_approved", true)
			if err != nil {
				response := models.Reply{
					Error:   err.Error(),
					Message: "error approving user",
					Success: true,
				}
				context.JSON(http.StatusOK, response)
				return
			}
			response := models.Reply{
				Data:    users.User{},
				Message: "succesfuly approved the user",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		}
	}
}
func RevokeUser(context *gin.Context) {
	id := context.Query("id")

	id = strings.ReplaceAll(id, "'", "")

	currentAdmin, err := CurrentUser(context)

	if err != nil {
		globalutils.UnAuthenticated(context)
		return
	} else if currentAdmin.AdminName == "" {
		globalutils.UnAuthorized(context)
		return
	} else {

		userExists, err := users.FindUserById(strings.ReplaceAll(id, "'", ""))
		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "error finding user",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if userExists.Firstname == "" {
			response := models.Reply{
				Error:   errors.New("user cannot be found").Error(),
				Message: "user does not exist",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if !userExists.IsApproved {
			response := models.Reply{
				Data:    users.User{},
				Message: "user is already revoked",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		} else {
			_, err := users.UpdateUserSpecificField(id, "is_approved", false)
			if err != nil {
				response := models.Reply{
					Error:   err.Error(),
					Message: "error revoking user",
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			}
			response := models.Reply{
				Data:    users.User{},
				Message: "succesfuly revoked the user",
				Success: true,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}
	}
}

func ApproveProduct(context *gin.Context) {
	productid := context.Query("id")

	id := strings.ReplaceAll(productid, "'", "")

	currentAdmin, err := CurrentUser(context)

	if err != nil {
		globalutils.UnAuthenticated(context)
		return
	} else if currentAdmin.AdminName == "" {
		globalutils.UnAuthorized(context)
		return
	} else {
		// check if product exist
		productExist, err := product.FindSingleProduct(id)
		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "error finding product",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if productExist.ProductName == "" {
			response := models.Reply{
				Message: "the product does not exist",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return

		} else if productExist.IsDeleted {
			response := models.Reply{
				Data:    productExist,
				Message: "cannot approve a deleted product!!Please restore product first",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		} else if productExist.IsApproved {
			response := models.Reply{
				Data:    productExist,
				Message: "product is already approved",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		} else {

			success, err := ApproveAd(id)
			fmt.Printf("this is query \n%s\n", id)
			if err != nil {
				response := models.Reply{
					Error:   err.Error(),
					Message: "error approving  product",
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			} else if !success {
				response := models.Reply{
					Error:   errors.New("could not approve product!!try again").Error(),
					Message: "failed in approving product",
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			} else {
				response := models.Reply{
					Data:    productExist,
					Message: "succesfully approved the product",
					Success: true,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			}
		}
	}

}
