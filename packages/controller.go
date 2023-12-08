package packages

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePackage(context *gin.Context) {
	var packageInput AddPackage

	if err := context.ShouldBind(&packageInput); err != nil {
		response := models.Reply{
			Message: "could not bind data from the user",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	packageuuid := uuid.New()
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// check if category exists
	packageExist, err := QueryPackageByName(strings.ReplaceAll(packageInput.PackageName, "'", ""))
	if err != nil {
		response := models.Reply{
			Message: "error fetching package",
			Success: false,
			Error:   errors.New("error fetching package").Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if packageExist.PackageName != "" {
		response := models.Reply{
			Message: "Package with similar name exists",
			Success: false,
			Error:   errors.New("package exists").Error(),
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		if packageInput.Duration != 0 && packageInput.PackageName != "" {
			newPackage := PackageModel{
				PackageId:   packageuuid.String(),
				PackageName: packageInput.PackageName,
				UsersNumber: 0,
				Price:       packageInput.Price,
				Duration:    packageInput.Duration,
				DateCreated: currentTime,
				DateUpdated: currentTime,
			}
			savedPackage, err := newPackage.Save()
			if err != nil {
				response := models.Reply{
					Message: "could not save package",
					Success: false,
					Error:   err.Error(),
				}
				context.JSON(http.StatusInternalServerError, response)
				return
			}
			response := models.Reply{
				Message: "package has been added succesfully",
				Success: true,
				Data:    savedPackage,
			}
			context.JSON(http.StatusCreated, response)
			return
		} else {
			response := models.Reply{
				Message: "contains null values",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}
	}

}

func UpdatePackage(context *gin.Context) {

	var packageInput PackageModel

	if err := context.ShouldBindJSON(&packageInput); err != nil {
		response := models.Reply{
			Message: "could not bind the user data to the request needs",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		packageid := strings.ReplaceAll(context.Query("id"), "'", "")

		packageExist, err := QuerySinglePackageUtil(strings.ReplaceAll(packageid, " ", ""))

		if err != nil {
			response := models.Reply{
				Message: " error querying the package",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if packageExist.PackageName == "" {
			response := models.Reply{
				Message: "package does not exist",
				Success: false,
			}
			context.JSON(http.StatusOK, response)
			return
		}
		packageExistName, err := QueryPackageByName(strings.ReplaceAll(packageInput.PackageName, "'", ""))
		if err != nil {
			response := models.Reply{
				Message: "error fetching package",
				Success: false,
				Error:   errors.New("error fetching package").Error(),
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if packageExistName.PackageName != "" {
			response := models.Reply{
				Message: "Package with such name does exists",
				Success: false,
				Error:   errors.New("package does exists").Error(),
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}

		currentTime := time.Now().Format("2006-01-02 15:04:05")
		id := strings.ReplaceAll(packageExist.PackageId, "'", "")

		newPackage := PackageModel{
			PackageName: packageInput.PackageName,
			UsersNumber: packageInput.UsersNumber,
			Price:       packageInput.Price,
			Duration:    packageInput.Duration,
			DateUpdated: currentTime,
		}
		packageUpdate, err := UpdatePackageUtil(id, newPackage)

		if err != nil {
			response := models.Reply{
				Message: "could not update package",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if packageUpdate.PackageName == "" {
			response := models.Reply{
				Error:   errors.New("could not update package").Error(),
				Message: "could not update package",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			response := models.Reply{
				Message: "Package updated",
				Success: true,
				Data:    packageUpdate,
			}
			context.JSON(http.StatusOK, response)
			return
		}
	}
}

func FetchSinglePackage(context *gin.Context) {
	packageid := strings.ReplaceAll(context.Query("id"), "'", "")

	packageExist, err := QuerySinglePackageUtil(packageid)
	if err != nil {
		response := models.Reply{
			Message: "could not fetch package",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if packageExist.PackageName == "" {
		response := models.Reply{
			Message: "package does not exist",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		response := models.Reply{
			Message: "package fetched",
			Success: true,
			Data:    packageExist,
		}
		context.JSON(http.StatusOK, response)
		return
	}
}

func FetchAllPackages(context *gin.Context) {
	packages, err := Fetchproducts()
	if err != nil {
		response := models.Reply{
			Message: "error fetching packages",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if len(packages) < 1 {
		response := models.Reply{
			Message: "There are no packages",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		response := models.Reply{
			Message: "Packages fetched successfuly",
			Success: true,
			Data:    packages,
		}
		context.JSON(http.StatusOK, response)
		return
	}
}
