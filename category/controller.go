package category

import (
	"net/http"

	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCategory(context *gin.Context) {
	var categoryInput models.Category

	if err := context.ShouldBindJSON(&categoryInput); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"category_input_error": err.Error(),
			"success":              false,
			"message":              "Wrong input from user",
		})
		return
	} else {
		// check if category already exists
		category, err := FetchSingleCategory(categoryInput.CategoryName)
		if err != nil {
			response := models.Reply{
				Message: "error validating the request",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if category.CategoryName != "" {
			response := models.Reply{
				Message: "category already exists",
				Data:    category,
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			categoryuuid := uuid.New()

			newProduct := models.Category{
				CategoryID:    categoryuuid.String(),
				CategoryName:  categoryInput.CategoryName,
				CategoryImage: categoryInput.CategoryImage,
			}

			category, err := newProduct.Save()

			if err != nil {
				response := models.Reply{
					Error:   err.Error(),
					Message: "Could not create category",
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			}

			response := models.Reply{
				Message: "category created succesfuly",
				Success: true,
				Data:    category,
			}
			context.JSON(http.StatusCreated, response)
		}

	}

}

func GetCategories(context *gin.Context) {
	categories, err := FetchAllCategories()
	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error fetching categories",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
	} else {
		response := models.Reply{
			Message: "fetched categories succesful",
			Success: true,
			Data:    categories,
		}
		context.JSON(http.StatusOK, response)
	}
}

func DeleteCategory(context *gin.Context) {
	categoryname := context.Param("name")

	// check if category exists
	category, err := FetchSingleCategory(categoryname)
	if err != nil {
		response := models.Reply{
			Message: "error checking the validity of query",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if category.CategoryName == "" {
		response := models.Reply{
			Message: "the category requested is missing!!please confirm the validity of the request",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if category.IsDeleted {
		response := models.Reply{
			Message: "the category requested is already deleted!!please confirm the validity of the request",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		response := models.Reply{
			Message: "delete operation succesful!!category deleted",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	}

}
