package mainad

import (
	"errors"
	"regexp"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
)

func ValidateMainAdInput(mainad *models.MainAd) (bool, error) {
	mainaddetails := []string{mainad.AdImage, mainad.AdCategory}
	charPattern := "[!@#$%^&*()_+\\-=\\[\\]{};:\\\\|,.<>?]"
	for _, value := range mainaddetails {
		if value == mainad.AdBy {
			if len(value) < 3 {
				return false, errors.New("ad owner id should be longer than 3 character long")
			}
		} else if value == mainad.AdImage {
			if len(value) < 3 {
				return false, errors.New("ad image string should be longer than 3 character long")
			}
		} else if value == mainad.Advertid {
			if len(value) < 3 {
				return false, errors.New("ad id should be longer than 3 character long")
			}
		} else if value == mainad.AdCategory {
			if len(value) < 3 {
				return false, errors.New("ad category should atleast be 3 characters long")
			} else if regexp.MustCompile(charPattern).MatchString(mainad.AdCategory) {
				return false, errors.New("ad category should not contain special character")
			}
		}

	}
	return true, nil
}

func GetAllMainAdsUtil() ([]models.MainAd, error) {
	var allmainads []models.MainAd
	err := database.Database.Find(&allmainads).Error
	if err != nil {
		return []models.MainAd{}, err
	}

	return allmainads, nil
}
func GetSingleMainAdUtil(query string) (models.MainAd, error) {
	var singlemainad models.MainAd
	err := database.Database.Where(query).Find(&singlemainad).Error
	if err != nil {
		return models.MainAd{}, err
	}
	return singlemainad, nil
}

func UpdateMainAdutil(query string, update models.MainAd) (models.MainAd, error) {
	var updatedAd models.MainAd

	result := database.Database.Model(&updatedAd).Where(query).Updates(update)
	if result.RowsAffected == 0 {
		return models.MainAd{}, errors.New("could not update the main ad")
	}
	return updatedAd, nil
}
func RestoreAdUtil(query string) (models.MainAd, error) {
	var updatedAd models.MainAd

	result := database.Database.Model(&updatedAd).Where(query).Update("is_deleted", false)
	if result.RowsAffected == 0 {
		return models.MainAd{}, errors.New("could not update the main ad")
	}
	return updatedAd, nil
}

func DeactivateUtil(query string) (models.MainAd, error) {
	var updatedAd models.MainAd

	result := database.Database.Model(&updatedAd).Where(query).Update("is_active", false)
	if result.RowsAffected == 0 {
		return models.MainAd{}, errors.New("could not update the main ad")
	}
	return updatedAd, nil
}
