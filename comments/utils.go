package comments

import (
	"errors"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
)

type Commentinput struct {
	ProductID string `json:"productid"`
	Comment   string `json:"comment"`
}

func GetProductComments(productid string) ([]models.Comment, error) {
	var comments []models.Comment
	err := database.Database.Where("isdeleted=?", false).Where("Product_ID=?", productid).Find(&comments).Error
	if err != nil {
		return []models.Comment{}, err

	}

	return comments, nil
}

func DeleteCommentUtil(commentId string, update models.Comment) (models.Comment, error) {
	var deletedComment models.Comment
	result := database.Database.Model(deletedComment).Where("comment_id=?", commentId).Updates(update)
	if result.RowsAffected == 0 {
		return models.Comment{}, errors.New("could not delete the comment")
	}
	return deletedComment, nil
}

func FetchComment(commentid string) (models.Comment, error) {
	var commentExists models.Comment
	err := database.Database.Where("comment_id=?", commentid).Find(&commentExists).Error
	if err != nil {
		return models.Comment{}, err
	}
	return commentExists, nil
}
