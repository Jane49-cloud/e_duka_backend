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

	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func GenerateJWT(user models.User) (string, error) {
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

func CurrentUser(context *gin.Context) (models.User, error) {
	err := ValidateJWT(context)
	if err != nil {
		return models.User{}, err
	}
	token, _ := GetToken(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	useremail := string(claims["email"].(string))

	user, err := FindUserByEmail(useremail)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func ValidateRegisterInput(user *RegisterInput) (bool, error) {
	userDetails := []string{user.Email, user.Firstname, user.Lastname, user.UserLocation, user.UserImage, user.Phone, user.Password}
	charPattern := "[!@#$%^&*()_+\\-=\\[\\]{};':\"\\\\|,.<>?]"
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
	fmt.Printf("validating login input\n%v", user)
	charPattern := "[!@#$%^&*()_+\\-=\\[\\]{};':\"\\\\|,.<>?]"
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
