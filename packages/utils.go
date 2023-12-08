package packages

import (
	"errors"
	"fmt"

	"eleliafrika.com/backend/database"
)

func QuerySinglePackageUtil(id string) (PackageModel, error) {
	var packageModel PackageModel
	err := database.Database.Where("package_id=?", id).Find(&packageModel).Error
	if err != nil {
		return PackageModel{}, err

	}
	fmt.Println(packageModel)
	return packageModel, nil
}
func QueryPackageByName(name string) (PackageModel, error) {
	var packageModel PackageModel
	err := database.Database.Where("package_name=?", name).Find(&packageModel).Error
	if err != nil {
		return PackageModel{}, err

	}
	return packageModel, nil
}

func Fetchproducts() ([]PackageModel, error) {
	var packagesList []PackageModel

	err := database.Database.Find(&packagesList).Error
	if err != nil {
		return []PackageModel{}, err
	}
	return packagesList, nil
}
func UpdatePackageUtil(id string, update PackageModel) (PackageModel, error) {
	var updatedPackage PackageModel

	result := database.Database.Model(&updatedPackage).Where("package_id=?", id).Updates(update)

	if result.RowsAffected == 0 {
		return PackageModel{}, errors.New("could not update the product right now")
	}
	return updatedPackage, nil
}
