package comments

import (
	"net/http"
	"time"

	product "eleliafrika.com/backend/Product"
	"eleliafrika.com/backend/models"
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MakeComment(context *gin.Context) {
	var commentinput Commentinput

	if err := context.ShouldBindJSON(&commentinput); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"success": false,
			"message": "could not make comments",
		})
	}
	success, err := ValidateCommentInput(&commentinput)
	if err != nil {
		response := models.Reply{
			Message: "Error validating user input",
			Error:   err.Error(),
			Success: false,
			Data:    commentinput,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if !success {
		response := models.Reply{
			Message: "Error validating user input",
			Error:   err.Error(),
			Success: false,
			Data:    commentinput,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {

		commentuuid := uuid.New()
		currentTime := time.Now()
		formattedTime := currentTime.Format("2006-01-02 15:04:05")
		user, err := users.CurrentUser(context)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "could not find user",
			})
			return
		}
		comment := models.Comment{
			CommentID:     commentuuid.String(),
			ProductID:     commentinput.ProductID,
			UserID:        user.UserID,
			Comment:       commentinput.Comment,
			DateCommented: formattedTime,
		}

		commentMade, err := comment.Save()

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"success": false,
				"message": "Error saving the comment",
			})
		} else {
			context.JSON(http.StatusCreated, gin.H{
				"success": true,
				"message": "Comment made",
				"comment": commentMade,
			})

		}

	}
}
func GetComments(context *gin.Context) {
	productid := context.Param("id")
	comments, err := GetProductComments(productid)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success":        false,
			"message":        "could not fetch comments for the products",
			"comments error": err.Error(),
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  "Succesfully fetched comments",
			"comments": comments,
		})
	}

}
func DeleteComment(context *gin.Context) {
	var commentinput models.Commentinput

	if err := context.ShouldBindJSON(&commentinput); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"err":     err.Error(),
			"message": "Error could not delete the comment due to bad body",
		})
	}

	// check if product exists and is visible to the users
	productExists, err := product.FindSingleProduct(commentinput.ProductID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err,
			"message": "an error occurred while fetching the product related to the comment",
		})
	} else if productExists.ProductID == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "The product is either inactive, deleted or simply does not exist",
		})

	} else {

		// check if the comment exist before deleting comment
		commentExists, err := FetchComment(commentinput.CommentID)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err,
				"message": "an error occurred",
			})
			return
		} else if commentExists.CommentID == "" {
			context.JSON(http.StatusOK, gin.H{
				"success": "false",
				"message": "Could not find the comment",
			})
			return
		} else {
			if commentExists.Isdeleted {
				context.JSON(http.StatusOK, gin.H{
					"message": "Comment already deleted",
					"success": false,
				})
			} else {
				commentDeleted, err := DeleteCommentUtil(commentinput.CommentID, models.Comment{
					Isdeleted: true,
				})

				if err != nil {
					context.JSON(http.StatusBadRequest, gin.H{
						"message": "Could not delete comment",
						"comment": commentDeleted,
						"id":      commentinput.CommentID,
						"error":   err.Error(),
					})
				} else {
					context.JSON(http.StatusOK, gin.H{
						"exist":   commentExists,
						"message": "Comment succesfully deleted",
					})
				}
			}

		}
	}

}
