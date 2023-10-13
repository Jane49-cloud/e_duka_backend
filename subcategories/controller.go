package subcategory

import (
	"net/http"

	"eleliafrika.com/backend/category"
	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateSubCategory(context *gin.Context) {
	var subcategoryInput models.SubCategory

	if err := context.ShouldBindJSON(&subcategoryInput); err != nil {
		response := models.Reply{
			Message: "Wrong input from the user",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
	success, err := ValidateSubCategoryInput(&subcategoryInput)
	if err != nil {
		response := models.Reply{
			Message: "error validating user input",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if !success {
		response := models.Reply{
			Message: "error validating user input for sub category",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		// check if the parent category exists
		parentcategory, err := category.FetchSingleCategory(subcategoryInput.ParentCategory)

		if err != nil {
			response := models.Reply{
				Message: "error validating the request for parent",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			if parentcategory.CategoryName == "" {
				response := models.Reply{
					Message: "parent category not found!!Please validate the data being passed",
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			} else {

				// check if sub category already exists
				subcategory, err := FetchSingleSubCategory(subcategoryInput.SubCategoryName)
				if err != nil {
					response := models.Reply{
						Message: "error validating the request",
						Error:   err.Error(),
						Success: false,
						Data:    subcategory,
					}
					context.JSON(http.StatusBadRequest, response)
					return
				} else if subcategory.SubCategoryName != "" {

					response := models.Reply{
						Message: "sub category already exists",
						Data:    subcategory,
						Success: false,
					}
					context.JSON(http.StatusBadRequest, response)
					return
				} else {
					categoryuuid := uuid.New()

					newSubCategory := models.SubCategory{
						SubCategoryID:    categoryuuid.String(),
						SubCategoryName:  subcategoryInput.SubCategoryName,
						SubCategoryImage: subcategoryInput.SubCategoryImage,
						ParentCategory:   subcategoryInput.ParentCategory,
					}

					category, err := newSubCategory.Save()

					if err != nil {
						response := models.Reply{
							Error:   err.Error(),
							Message: "Could not create sub category",
							Success: false,
						}
						context.JSON(http.StatusBadRequest, response)
						return
					}

					response := models.Reply{
						Message: "sub category created succesfuly",
						Success: true,
						Data:    category,
					}
					context.JSON(http.StatusCreated, response)
				}
			}
		}

	}

}

func GetSubCategories(context *gin.Context) {
	categoryname := context.Param("name")
	subCategories, err := FetchAllSubCategories(categoryname)
	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error fetching sub categories",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
	} else {
		response := models.Reply{
			Message: "fetched sub categories succesful",
			Success: true,
			Data:    subCategories,
		}
		context.JSON(http.StatusOK, response)
	}
}

func DeleteSubCategory(context *gin.Context) {
	subcategoryname := context.Param("name")

	// check if category exists
	subCategoryExist, err := FetchSingleSubCategory(subcategoryname)
	if err != nil {
		response := models.Reply{
			Message: "error checking the validity of query",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if subCategoryExist.SubCategoryName == "" {
		response := models.Reply{
			Message: "the sub category requested is missing!!please confirm the validity of the request",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if subCategoryExist.IsDeleted {
		response := models.Reply{
			Message: "the sub category requested is already deleted!!please confirm the validity of the request",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		deletedCategory, err := UpdateSubCategory(subcategoryname, models.SubCategory{
			IsDeleted: true,
		})
		if err != nil {
			response := models.Reply{
				Message: "Could not delete the sub category",
				Success: false,
				Error:   err.Error(),
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {

			response := models.Reply{
				Message: "delete operation succesful!!sub category deleted",
				Success: true,
				Data:    deletedCategory,
			}
			context.JSON(http.StatusOK, response)
			return
		}
	}

}
