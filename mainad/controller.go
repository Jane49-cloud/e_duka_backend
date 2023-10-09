package mainad

import (
	"net/http"
	"time"

	"eleliafrika.com/backend/category"
	"eleliafrika.com/backend/models"
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateMainAd(context *gin.Context) {
	var createAdInput models.MainAd

	if err := context.ShouldBindJSON(&createAdInput); err != nil {
		response := models.Reply{
			Message: "could not bind data from the user",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	adid := uuid.New()
	currentTime := time.Now()
	var currentuserId string
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	success, err := ValidateMainAdInput(&createAdInput)

	if err != nil {
		response := models.Reply{
			Message: "could not validate the data of the created ad",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if !success {
		response := models.Reply{
			Message: "verifying data was not succesful",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		currentuser, err := users.CurrentUser(context)
		if err != nil {
			response := models.Reply{
				Message: "could not find current user",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if currentuser.Firstname == "" {
			response := models.Reply{
				Message: "user not found",
				Success: false,
			}
			context.JSON(http.StatusNetworkAuthenticationRequired, response)
			return
		} else {
			currentuserId = currentuser.UserID
			newMainAd := models.MainAd{
				Advertid:    adid.String(),
				AdBy:        currentuserId,
				AdImage:     createAdInput.AdImage,
				AdName:      createAdInput.AdName,
				DateCreated: formattedTime,
				AdCategory:  createAdInput.AdCategory,
			}

			// check if the caegory exists
			categoryExists, err := category.FetchSingleCategory(createAdInput.AdCategory)
			if err != nil {
				response := models.Reply{
					Message: "error fetching the category provided",
					Error:   err.Error(),
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			} else if categoryExists.CategoryName == "" {
				response := models.Reply{
					Message: "the category provided does not exist",
					Success: false,
					Data:    newMainAd,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			}
			ad, err := newMainAd.Save()
			if err != nil {
				response := models.Reply{
					Message: "error creating the ad",
					Error:   err.Error(),
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			} else {
				response := models.Reply{
					Message: "ad created succesfully",
					Success: true,
					Data:    ad,
				}
				context.JSON(http.StatusOK, response)
				return
			}
		}

	}

}
func GetAllMainAds(context *gin.Context) {

}
func UpdateMainAd(context *gin.Context) {

}
func GetSingleMainAd(context *gin.Context) {

}
