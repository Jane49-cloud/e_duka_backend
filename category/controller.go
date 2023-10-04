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
	}

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
