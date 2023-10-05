package brands

import (
	"net/http"

	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddBrand(context *gin.Context) {
	var brandInput models.Brand

	if err := context.ShouldBindJSON(&brandInput); err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "Could not bind the data from user",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
	}
	success, err := ValidatebrandInput(&brandInput)
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
			Message: "error validating user input for brand",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {

		// check if brand already exists
		brandExist, err := FetchSingleBrand(brandInput.BrandName)
		if err != nil {
			response := models.Reply{
				Message: "could not validate request",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if brandExist.BrandID != "" {
			response := models.Reply{
				Message: "Brand name already exist",
				Success: true,
				Data:    brandExist,
			}
			context.JSON(http.StatusOK, response)
			return
		} else {

			branduuid := uuid.New()
			newBrand := models.Brand{
				BrandID:          branduuid.String(),
				BrandName:        brandInput.BrandName,
				Imageurl:         brandInput.Imageurl,
				TotalProducts:    0,
				TotalEngagements: 0,
			}
			brand, err := newBrand.Save()
			if err != nil {
				response := models.Reply{
					Message: "error occurred during creation",
					Error:   err.Error(),
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
			} else {
				response := models.Reply{
					Message: "success adding brands",
					Success: true,
					Data:    brand,
				}
				context.JSON(http.StatusOK, response)

			}
		}
	}
}

func GetAllBrands(context *gin.Context) {
	brands, err := FetchAllBrands()
	if err != nil {
		response := models.Reply{
			Message: "Error fetching brands",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
	} else {
		response := models.Reply{
			Message: "fetched all brands",
			Success: true,
			Data:    brands,
		}
		context.JSON(http.StatusOK, response)
	}
}
func DeleteBrand(context *gin.Context) {
	brandname := context.Param("name")
	// check if brand exist
	brand, err := FetchSingleBrand(brandname)
	if err != nil {
		response := models.Reply{
			Message: "Could not validate request",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if brand.BrandName == "" {
		response := models.Reply{
			Message: "brand does not exist in the database",
			Success: false,
			Data:    brand,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		brandDeleted, err := UpdateBrand(brandname, models.Brand{
			Isdeleted: true,
		})

		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "Error occurred deleting the brand",
				Success: false,
				Data:    brandDeleted,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			response := models.Reply{
				Message: "Brand deleted successfully",
				Success: true,
				Data:    brandDeleted,
			}
			context.JSON(http.StatusOK, response)
			return
		}
	}
}
