package admin

import (
	"errors"
	"fmt"
	"strings"
	"time"

	globalutils "eleliafrika.com/backend/global_utils"
	"eleliafrika.com/backend/product"
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(context *gin.Context) {
	var input AddAdmin

	if err := context.ShouldBindJSON(&input); err != nil {

		globalutils.HandleError("could not bind json data from user", err, context)
		return
	}
	success, err := ValidateRegisterInput(&input)
	if err != nil {

		globalutils.HandleError("error validating user details", err, context)
		return
	}
	if !success {

		globalutils.HandleError("returned false in validating user data", errors.New("could validate data"), context)
		return
	}

	// get current data to save user with
	randomuuid := uuid.New()
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	userImagepath, err := users.UploadUserImage(input.AdminImage, input.AdminName)

	if err != nil {
		globalutils.HandleError("error uploading user image", err, context)
		return
	}

	admin := SystemAdmin{
		AdminID:      randomuuid.String(),
		AdminName:    input.AdminName,
		Email:        input.Email,
		AdminImage:   userImagepath,
		Cell:         input.Cell,
		Password:     input.Password,
		DateAdded:    formattedTime,
		LastLoggedIn: formattedTime,
	}

	userExists, err := FindAdminByEmail(admin.Email)

	if err != nil {
		globalutils.HandleError("error validating user", err, context)
		return
	} else if userExists.Email != "" {
		globalutils.HandleError("email already used", errors.New("user email has already been used"), context)
		return
	} else {

		_, err := admin.Save()
		if err != nil {
			globalutils.HandleError("error creating user", err, context)
			return
		}
		// generate token directly on succesfuly register
		token, err := GenerateJWT(admin)

		if err != nil {
			globalutils.HandleError("error generatin token for user", err, context)
			return
		}

		globalutils.HandleSuccess("admin added", token, context)
	}
}

func Login(context *gin.Context) {
	var input AdminLogin

	if err := context.ShouldBindJSON(&input); err != nil {

		globalutils.HandleError("errror binding data json", err, context)
		return
	}
	// check validity of user input
	success, err := ValidateLoginInput(&input)
	if err != nil {

		globalutils.HandleError("error validating user input", err, context)
		return
	} else if !success {
		globalutils.HandleError("error validating user", err, context)
		return
	} else {

		// check if user exists

		admin, err := FindAdminByEmail(input.Email)

		if err != nil {
			globalutils.HandleError("error fetching user", err, context)
			return
		}

		// validate the password password passed with the harsh on db

		err = admin.ValidatePassword(input.Password)
		if err != nil {
			globalutils.HandleError("incorect password: "+input.Password, err, context)
			return
		}

		// generate jwt if error does not exists
		token, err := GenerateJWT(admin)

		if err != nil {

			globalutils.HandleError("error while generating token", err, context)
			return

		}

		globalutils.HandleSuccess("login succesfull", token, context)
		return
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
			globalutils.HandleError("error fetching all users", err, context)
			return
		} else {
			globalutils.HandleSuccess("succesfully fetched all users", users, context)
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
			globalutils.HandleError("error finding user", err, context)
			return
		} else if userExists.Firstname == "" {
			globalutils.HandleError("user does not exist", errors.New("user cannot be found"), context)
			return
		} else if userExists.IsApproved {
			globalutils.HandleSuccess("user is already approved", users.User{}, context)
			return
		} else {

			fmt.Printf("request id\n%v\n", id)
			_, err := users.UpdateUserSpecificField(id, "is_approved", true)
			if err != nil {
				globalutils.HandleError("error approving user", err, context)
				return
			}
			globalutils.HandleSuccess("succesfuly approved the user", users.User{}, context)
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
			globalutils.HandleError("error finding user", err, context)
			return
		} else if userExists.Firstname == "" {
			globalutils.HandleError("user does not exist", errors.New("user cannot be found"), context)
			return
		} else if !userExists.IsApproved {
			globalutils.HandleSuccess("user is already revoked", users.User{}, context)
			return
		} else {
			_, err := users.UpdateUserSpecificField(id, "is_approved", false)
			if err != nil {
				globalutils.HandleError("error revoking user", err, context)
				return
			}
			globalutils.HandleSuccess("succesfuly revoked the user", users.User{}, context)
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
			globalutils.HandleError("error finding product", err, context)
			return
		} else if productExist.ProductName == "" {
			globalutils.HandleSuccess("the product does not exist", product.Product{}, context)
			return
		} else if productExist.IsDeleted {
			globalutils.HandleSuccess("cannot approve a deleted product!!Please restore product first", productExist, context)
			return
		} else if productExist.IsApproved {
			globalutils.HandleSuccess("product is already approved", productExist, context)
			return
		} else {

			success, err := ApproveAd(id)
			fmt.Printf("this is query \n%s\n", id)
			if err != nil {
				globalutils.HandleError("error approving  product", err, context)
				return
			} else if !success {
				globalutils.HandleError("failed in approving product", errors.New("could not approve product!!try again"), context)
				return
			} else {
				globalutils.HandleSuccess("succesfully approved the product", productExist, context)
				return
			}
		}
	}

}