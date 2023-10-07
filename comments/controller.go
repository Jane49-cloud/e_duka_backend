package comments

import (
	"fmt"
	"net/http"
	"time"

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
	commentId := context.Param("id")
	query := "comment_id=" + commentId

	// check if the comment exist before deleting comment
	commentExists, err := FetchComment(query)
	fmt.Printf("this is comment\n%v\n", commentExists)
	if err != nil {
		response := models.Reply{
			Message: "an error occured during fetching comment",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if commentExists.CommentID == "" {
		response := models.Reply{
			Message: "comment does not exist",
			Success: false,
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		if commentExists.Isdeleted {
			response := models.Reply{
				Message: "comment already deleted",
				Success: false,
			}
			context.JSON(http.StatusOK, response)
		} else {
			fmt.Printf("id\n%s\n", commentId)
			commentDeleted, err := DeleteCommentUtil(query, models.Comment{
				Isdeleted: true,
			})

			if err != nil {
				response := models.Reply{
					Message: "Could not delete comment",
					Error:   err.Error(),
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
			} else {
				response := models.Reply{
					Message: "comment deleted succesfully",
					Data:    commentDeleted,
					Success: true,
				}
				context.JSON(http.StatusOK, response)
			}
		}

	}
}
