package users

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"eleliafrika.com/backend/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func GenerateJWT(user User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"eat":   time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}

func ValidateJWT(context *gin.Context) error {
	token, err := GetToken(context)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

func GetToken(context *gin.Context) (*jwt.Token, error) {
	tokenFromUser := context.Request.Header.Get("x-access-token")

	token, err := jwt.Parse(tokenFromUser, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return privateKey, nil
	})

	return token, err
}

func JWTAuthMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			context.Abort()
			return
		}
		context.Next()
	}
}

func CurrentUser(context *gin.Context) (User, error) {
	err := ValidateJWT(context)
	if err != nil {
		return User{}, err
	}
	token, _ := GetToken(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	useremail := string(claims["email"].(string))

	user, err := FindUserByEmail(useremail)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func ValidateRegisterInput(user *RegisterInput) (bool, error) {
	userDetails := []string{user.Email, user.Firstname, user.Lastname, user.UserLocation, user.Phone, user.Password}
	charPattern := "[!@#$%^&*()_+\\-=\\[\\]{};':\"\\\\|,.<>?]"
	numPattern := "[0-9]"
	capPattern := "[A-Z]"
	for _, value := range userDetails {
		value = strings.ReplaceAll(value, " ", "")
		if value == user.Email {
			user.Email = strings.ToLower(user.Email)
			emailPattern := "[!#$%^&*()+\\=\\[\\]{};':\"\\\\|,<>?]"
			if len(value) < 8 {
				return false, errors.New("invalid email format!email should contain @, . and should be longer than 8 characters")
			} else if !strings.Contains(user.Email, "@") {
				return false, errors.New("invalid email format!email should contain @, . and should be longer than 8 characters")
			} else if !strings.Contains(user.Email, ".") {
				return false, errors.New("invalid email format!email should contain @, . and should be longer than 8 characters")
			} else if regexp.MustCompile(emailPattern).MatchString(user.Email) {
				return false, errors.New("invalid characters in email")
			}
		} else if value == user.Phone {
			if len(value) < 10 {
				return false, errors.New("phone number should be atleast 10 characters long")
			} else {
				for _, char := range user.Phone {
					if !unicode.IsNumber(char) {
						return false, errors.New("phone number can only contain numbers")
					}
				}
			}
		} else if value == user.Password {

			if len(value) < 8 {
				return false, errors.New("password is too short")
			} else if !regexp.MustCompile(charPattern).MatchString(user.Password) {
				return false, errors.New("password must contain atleast one special character")
			} else if !regexp.MustCompile(numPattern).MatchString(user.Password) {
				return false, errors.New("password must contain atleast numerical digit")
			} else if !regexp.MustCompile(capPattern).MatchString(user.Password) {
				return false, errors.New("password must contain a capital letter")
			}
		} else if value == user.UserLocation {
			user.UserLocation = strings.ToLower(user.UserLocation)
			if len(value) < 3 {
				return false, errors.New("location is too short")
			} else if regexp.MustCompile(numPattern).MatchString(user.UserLocation) {
				return false, errors.New("location must not contain a numerical digit")
			}
		} else if len(value) < 3 {
			return false, errors.New("invalid input length for field")
		} else {
			if regexp.MustCompile(numPattern).MatchString(user.Firstname) {
				return false, errors.New("first name cannot contain numbers")
			} else if regexp.MustCompile(charPattern).MatchString(user.Firstname) {
				return false, errors.New("first name cannot contain special characters")
			}
			if regexp.MustCompile(numPattern).MatchString(user.Middlename) {
				return false, errors.New("middle name cannot contain numbers")
			} else if regexp.MustCompile(charPattern).MatchString(user.Middlename) {
				return false, errors.New("middle name cannot contain special characters")
			}
			if regexp.MustCompile(numPattern).MatchString(user.Lastname) {
				return false, errors.New("last name cannot contain numbers")
			} else if regexp.MustCompile(charPattern).MatchString(user.Lastname) {
				return false, errors.New("last name cannot contain special characters")
			}
		}
	}
	return true, nil
}

func ValidateLoginInput(user *LoginInput) (bool, error) {
	userDetails := []string{user.Email, user.Password}
	charPattern := "[!#$%^&*()+\\=\\[\\]{};':\"\\\\|,<>?@.-]"
	emailPattern := "[!#$%^&*()+\\=\\[\\]{};':\"\\\\|,<>?]"
	numPattern := "[0-9]"
	capPattern := "[A-Z]"
	for _, value := range userDetails {
		if value == user.Email {
			if len(value) < 8 {
				return false, errors.New("invalid email format!email should contain @, . and should longer than 8 characters")
			} else if !strings.Contains(user.Email, "@") {
				return false, errors.New("invalid email format!email should contain @, . and should longer than 8 characters")
			} else if !strings.Contains(user.Email, ".") {
				return false, errors.New("invalid email format!email should contain @, . and should longer than 8 characters")
			} else if regexp.MustCompile(emailPattern).MatchString(user.Email) {
				return false, errors.New("invalid characters in email")
			}
		} else if value == user.Password {
			if len(value) < 8 {
				return false, errors.New("password is too short")
			} else if !regexp.MustCompile(charPattern).MatchString(user.Password) {
				return false, errors.New("password must contain atleast one special character")
			} else if !regexp.MustCompile(numPattern).MatchString(user.Password) {
				return false, errors.New("password must contain atleast numerical digit")
			} else if !regexp.MustCompile(capPattern).MatchString(user.Password) {
				return false, errors.New("password must contain a capital letter")
			}
		}
	}
	return true, nil
}

func UpdateUserUtil(query string, update User) (User, error) {
	var updatedUser User

	result := database.Database.Model(&updatedUser).Where("user_id=?", query).Updates(update)
	if result.RowsAffected == 0 {
		return User{}, errors.New("could not update the user")
	}
	return updatedUser, nil
}
func UpdateUserSpecificField(userId string, field string, value any) (User, error) {
	var updatedUser User

	result := database.Database.Model(&updatedUser).Where("user_id=?", userId).Update(field, value)
	if result.RowsAffected == 0 {
		return User{}, errors.New("could not update")
	}
	return updatedUser, nil
}

func UploadUserImage(imageString string, username string) {

	// imageuuid := uuid.New()
	// filename := strings.ReplaceAll(username, " ", "") + imageuuid.String()

	// imageBytes, err := base64.StdEncoding.DecodeString(imageString)
	// if err != nil {
	// 	return "", err
	// }

	// imagepath, err := images.UploadImageToBucket(username, "user-images", imageBytes, filename)
	// if err != nil {
	// 	return "", nil
	// }
}

func FetchAllSellersUtil() ([]User, error) {
	var AllUsers []User

	err := database.Database.Where("is_approved=?", true).Find(&AllUsers).Error

	if err != nil {
		return []User{}, err
	}

	// if len(AllUsers) > 0 {
	// 	for _, user := range AllUsers {
	// 		userImage, err := images.DownloadImageFromBucket(user.UserImage)
	// 		if err != nil {
	// 			return []User{}, err
	// 		} else if userImage == "" {
	// 			return []User{}, errors.New("image not downloaded")
	// 		}
	// 		user.UserImage = userImage
	// 	}
	// }

	return AllUsers, nil
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
