package admin

import (
	"errors"
	"strings"

	globalutils "eleliafrika.com/backend/global_utils"
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func FetchSellers(context *gin.Context) {
	users, err := FetchAllUsersUtil()
	if err != nil {
		globalutils.HandleError("error fetching all users", err, context)
		return
	} else {
		globalutils.HandleSuccess("succesfully fetched all users", users, context)
		return
	}
}

func ApproveUser(context *gin.Context) {
	id := context.Query("id")

	query := "user_id=" + id

	userExists, err := users.FindUserById(strings.ReplaceAll(id, "'", ""))
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
		_, err := users.UpdateUserSpecificField(query, "is_approved", true)
		if err != nil {
			globalutils.HandleError("error approving user", err, context)
			return
		}
		globalutils.HandleSuccess("succesfuly approved the user", users.User{}, context)
	}
}
func RevokeUser(context *gin.Context) {
	id := context.Query("id")

	query := "user_id=" + id

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
		_, err := users.UpdateUserSpecificField(query, "is_approved", false)
		if err != nil {
			globalutils.HandleError("error revoking user", err, context)
			return
		}
		globalutils.HandleSuccess("succesfuly revoked the user", users.User{}, context)
	}
}
