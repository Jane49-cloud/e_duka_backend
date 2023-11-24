package admin

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/product"
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func ValidateRegisterInput(admin *AddAdmin) (bool, error) {
	details := []string{admin.AdminName, admin.Email, admin.Cell, admin.Password, admin.Role}

	charPattern := "[!@#$%^&*()_+\\-=\\[\\]{};':\"\\\\|,.<>?]"
	emailPattern := "[!#$%^&*()+\\=\\[\\]{};':\"\\\\|,<>?]"
	numPattern := "[0-9]"
	capPattern := "[A-Z]"
	for _, value := range details {
		if value == admin.AdminName {
			if len(value) < 5 {
				return false, errors.New("admin name is too short")
			} else if len(value) > 20 {
				return false, errors.New("admin name is too long")
			} else if regexp.MustCompile(charPattern).MatchString(admin.AdminName) {
				return false, errors.New("special characters not allowed in admin name")
			} else if regexp.MustCompile(numPattern).MatchString(admin.AdminName) {
				return false, errors.New("numerals not allowed in admin name")
			}
		} else if value == admin.Email {
			admin.Email = strings.ToLower(admin.Email)
			if len(value) < 8 {
				return false, errors.New("email should be longer than 8 characters")
			} else if !strings.Contains(admin.Email, "@") {
				return false, errors.New("email should contain @: " + value)
			} else if !strings.Contains(admin.Email, ".") {
				return false, errors.New("invalid email format!email should contain @, . and should be longer than 8 characters")
			} else if regexp.MustCompile(emailPattern).MatchString(admin.Email) {
				return false, errors.New("invalid characters in email")
			}
		} else if value == admin.Cell {
			if len(value) < 10 {
				return false, errors.New("phone number should be atleast 10 characters long")
			} else if len(value) > 15 {
				return false, errors.New("phone number is too long")
			} else {
				for _, char := range admin.Cell {
					if !unicode.IsNumber(char) {
						return false, errors.New("phone number can only contain numbers")
					}
				}
			}
		} else if value == admin.Password {
			if len(value) < 7 {
				return false, errors.New("password is too short")
			} else if !regexp.MustCompile(charPattern).MatchString(admin.Password) {
				return false, errors.New("password must contain atleast one special character")
			} else if !regexp.MustCompile(numPattern).MatchString(admin.Password) {
				return false, errors.New("password must contain atleast numerical digit")
			} else if !regexp.MustCompile(capPattern).MatchString(admin.Password) {
				return false, errors.New("password must contain a capital letter")
			}
		} else if value == admin.Role {
			admin.Role = strings.ToLower(admin.Role)
			if len(value) < 3 {
				return false, errors.New("role string is too short")
			} else if regexp.MustCompile(numPattern).MatchString(admin.Role) {
				return false, errors.New("admin role must not contain a numerical digit")
			}
		}
	}
	return true, nil
}
func ValidateLoginInput(admin *AdminLogin) (bool, error) {
	userDetails := []string{admin.Email, admin.Password}
	charPattern := "[!#$%^&*()+\\=\\[\\]{};':\"\\\\|,<>?@.-]"
	emailPattern := "[!#$%^&*()+\\=\\[\\]{};':\"\\\\|,<>?]"
	numPattern := "[0-9]"
	capPattern := "[A-Z]"
	for _, value := range userDetails {
		if value == admin.Email {
			if len(value) < 8 {
				return false, errors.New("invalid email format!email should contain @, . and should longer than 8 characters")
			} else if !strings.Contains(admin.Email, "@") {
				return false, errors.New("invalid email format!email should contain @, . and should longer than 8 characters")
			} else if !strings.Contains(admin.Email, ".") {
				return false, errors.New("invalid email format!email should contain @, . and should longer than 8 characters")
			} else if regexp.MustCompile(emailPattern).MatchString(admin.Email) {
				return false, errors.New("invalid characters in email")
			}
		} else if value == admin.Password {

			if len(value) < 8 {
				return false, errors.New("password is too short")
			} else if !regexp.MustCompile(charPattern).MatchString(admin.Password) {
				return false, errors.New("password must contain atleast one special character")
			} else if !regexp.MustCompile(numPattern).MatchString(admin.Password) {
				return false, errors.New("password must contain atleast numerical digit")
			} else if !regexp.MustCompile(capPattern).MatchString(admin.Password) {
				return false, errors.New("password must contain a capital letter")
			}
		}
	}
	return true, nil
}
func FetchAllUsersUtil() ([]users.User, error) {
	var AllUsers []users.User

	err := database.Database.Find(&AllUsers).Error

	if err != nil {
		return []users.User{}, err
	}

	return AllUsers, nil
}
func ApproveAd(id string) (bool, error) {
	var updatedProduct product.Product
	result := database.Database.Model(&updatedProduct).Where("product_id=?", id).Update("is_approved", true)
	if result.RowsAffected == 0 {
		return false, errors.New("could not approve the current product")
	}
	return true, nil
}
func FindAdminByEmail(email string) (SystemAdmin, error) {
	if len(email) < 10 {
		return SystemAdmin{}, errors.New("user email provided is null")
	}
	var admin SystemAdmin
	err := database.Database.Where("email=?", email).Find(&admin).Error
	if err != nil {
		return SystemAdmin{}, err
	}
	return admin, nil
}

func GenerateJWT(admin SystemAdmin) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": admin.Email,
		"iat":   time.Now().Unix(),
		"eat":   time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}
func (admin *SystemAdmin) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
}
func ValidateHashPassword(hashedPassword string, userPassword string) (bool, error) {
	if userPassword == "" {
		return false, bcrypt.ErrMismatchedHashAndPassword

	}

	if hashedPassword == "" {
		return false, bcrypt.ErrHashTooShort

	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))
	if err != nil {
		return false, err
	}
	return true, nil
}

func CurrentUser(context *gin.Context) (SystemAdmin, error) {
	err := users.ValidateJWT(context)
	if err != nil {
		return SystemAdmin{}, err
	}
	token, _ := users.GetToken(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	useremail := string(claims["email"].(string))

	user, err := FindAdminByEmail(useremail)
	if err != nil {
		return SystemAdmin{}, err
	}
	return user, nil
}
func UpdateAdminUtil(adminId string, update SystemAdmin) (SystemAdmin, error) {
	var updatedAdmin SystemAdmin
	result := database.Database.Model(&updatedAdmin).Where("admin_id=?", adminId).Updates(update)
	if result.RowsAffected == 0 {
		return SystemAdmin{}, errors.New("could not update")
	}
	return updatedAdmin, nil
}
