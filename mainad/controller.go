package mainad

import (
	"net/http"
	"time"

	"eleliafrika.com/backend/admin"
	"eleliafrika.com/backend/category"
	"eleliafrika.com/backend/models"
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
		currentuser, err := admin.CurrentUser(context)
		if err != nil {
			response := models.Reply{
				Message: "could not find current user",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if currentuser.AdminName == "" {
			response := models.Reply{
				Message: "user not found",
				Success: false,
			}
			context.JSON(http.StatusUnauthorized, response)
			return
		} else {
			currentuserId = currentuser.AdminID
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
	adlist, err := GetAllMainAdsUtil()
	if err != nil {
		response := models.Reply{
			Message: "error fetching the ads",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if len(adlist) == 0 {
		response := models.Reply{
			Message: "there are no ads",
			Success: true,
			Data:    adlist,
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		response := models.Reply{
			Message: "ads fetched",
			Success: true,
			Data:    adlist,
		}
		context.JSON(http.StatusOK, response)
		return
	}
}
func UpdateMainAd(context *gin.Context) {

}
func GetSingleMainAd(context *gin.Context) {
	adid := context.Query("id")
	query := "ad_id=" + adid

	singlead, err := GetSingleMainAdUtil(query)
	if err != nil {
		response := models.Reply{
			Message: "error fetching single ad",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if singlead.IsDeleted {
		response := models.Reply{
			Message: "the ad you are fetchin has been deleted",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else if !singlead.AdActive {
		response := models.Reply{
			Message: "the ad you are fetchin has been deactivated",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		response := models.Reply{
			Message: "succesffully fetched the ad",
			Success: true,
			Data:    singlead,
		}
		context.JSON(http.StatusOK, response)
		return
	}

}
func DeleteMainAd(context *gin.Context) {
	adid := context.Query("id")
	query := "ad_id=" + adid

	// check if ad exists and is not deleted or inactive
	adExist, err := GetSingleMainAdUtil(query)
	if err != nil {
		response := models.Reply{
			Message: "error fetching the ad",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if adExist.AdName == "" {
		response := models.Reply{
			Message: "ad does not exist",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else if adExist.IsDeleted {
		response := models.Reply{
			Message: "ad is already deleted",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else if adExist.AdActive {
		response := models.Reply{
			Message: "cannot delete an active ad!!First deactivate the ad to continue",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		_, err := UpdateMainAdutil(query, models.MainAd{
			IsDeleted: true,
		})
		if err != nil {
			response := models.Reply{
				Message: "error deleting the ad",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			response := models.Reply{
				Message: "ad has been deleted succesfully",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		}

	}
}

func RestoreMainAd(context *gin.Context) {
	adid := context.Query("id")
	query := "ad_id=" + adid

	// check if ad exists and is not deleted or inactive
	adExist, err := GetSingleMainAdUtil(query)
	if err != nil {
		response := models.Reply{
			Message: "error fetching the ad",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if adExist.AdName == "" {
		response := models.Reply{
			Message: "ad does not exist",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else if !adExist.IsDeleted {
		response := models.Reply{
			Message: "ad is not deleted",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		_, err := RestoreAdUtil(query)
		if err != nil {
			response := models.Reply{
				Message: "error restoring the ad",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			response := models.Reply{
				Message: "ad has been restored succesfully",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		}

	}
}

func ActivateMainAd(context *gin.Context) {
	adid := context.Query("id")
	query := "ad_id=" + adid

	// check if ad exists and is not deleted or inactive
	adExist, err := GetSingleMainAdUtil(query)
	if err != nil {
		response := models.Reply{
			Message: "error fetching the ad",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if adExist.AdName == "" {
		response := models.Reply{
			Message: "ad does not exist",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return

	} else if adExist.AdActive {
		response := models.Reply{
			Message: "ad is active",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else if adExist.IsDeleted {
		response := models.Reply{
			Message: "ad is deleted and operation in prohibitted",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		_, err := UpdateMainAdutil(query, models.MainAd{
			AdActive: true,
		})
		if err != nil {
			response := models.Reply{
				Message: "error activating the ad",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			response := models.Reply{
				Message: "ad has been activated succesfully",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		}

	}
}

func DeactivateMainAd(context *gin.Context) {
	adid := context.Query("id")
	query := "ad_id=" + adid

	// check if ad exists and is not deleted or inactive
	adExist, err := GetSingleMainAdUtil(query)
	if err != nil {
		response := models.Reply{
			Message: "error fetching the ad",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if adExist.AdName == "" {
		response := models.Reply{
			Message: "ad does not exist",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else if !adExist.AdActive {
		response := models.Reply{
			Message: "ad is not active",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		_, err := DeactivateUtil(query)
		if err != nil {
			response := models.Reply{
				Message: "error deactivating the ad",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			response := models.Reply{
				Message: "ad has been deactivated succesfully",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
		}

	}
}
