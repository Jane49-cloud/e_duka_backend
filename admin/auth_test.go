package admin

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidatePasssword(t *testing.T) {
	type SamplePassword struct {
		UserPassword   string
		Hashedpassword string
		Expected       bool
		Error          error
	}
	emptyPassword := SamplePassword{
		"",
		"$2a$10$nT6RqrRYwY3O.256u3NUPOUtTMQz08Qrl9VQT7daA3uIW.Ij3yAjm",
		false,
		bcrypt.ErrMismatchedHashAndPassword,
	}

	emptyHashedPassword := SamplePassword{
		"userpassword",
		"",
		false,
		bcrypt.ErrHashTooShort,
	}

	wrongUserPassword := SamplePassword{
		"userpassword",
		"$2a$10$nT6RqrRYwY3O.256u3NUPOUtTMQz08Qrl9VQT7daA3uIW.Ij3yAjm",
		false,
		bcrypt.ErrMismatchedHashAndPassword,
	}

	validPassword := SamplePassword{
		"@11Janejane",
		"$2a$10$nT6RqrRYwY3O.256u3NUPOUtTMQz08Qrl9VQT7daA3uIW.Ij3yAjm",
		true,
		nil,
	}

	cases := []SamplePassword{
		emptyPassword,
		emptyHashedPassword,
		wrongUserPassword,
		validPassword,
	}

	for _, password := range cases {
		valid, err := ValidateHashPassword(password.Hashedpassword, password.UserPassword)
		if valid != password.Expected {
			t.Errorf("test failed: expected %v but found %v error %v", password.Expected, valid, err)
		} else if password.Error != err {
			t.Errorf("test failed: expected %v but found %v", password.Error, err)
		}
	}
	t.Logf("test passed")
}

func TestRegAdminInput(t *testing.T) {
	type SampleAdmin struct {
		name     string
		data     AddAdmin
		expected bool
	}

	// cases for admin name
	emptyName := SampleAdmin{
		name:     "should return empty name",
		data:     AddAdmin{"", "user@gmail.com", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	blankName := SampleAdmin{
		name:     "should return blank name",
		data:     AddAdmin{"   ", "user@gmail.com", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	shortName := SampleAdmin{
		name:     "should return too short name",
		data:     AddAdmin{"as", "user@gmail.com", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	longName := SampleAdmin{
		name:     "should return too long name",
		data:     AddAdmin{"this is a very long name that is not allowed", "user@gmail.com", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	numInName := SampleAdmin{
		name:     "should return number are not allowed",
		data:     AddAdmin{"this1235", "user@gmail.com", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	specialCharInName := SampleAdmin{
		name:     "shoulname@#",
		data:     AddAdmin{"this is a very long name", "user@gmail.com", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}

	// cases for email
	emptyEmail := SampleAdmin{
		name:     "should return empty email",
		data:     AddAdmin{"this is a very long name", "", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	blankEmail := SampleAdmin{
		name:     "should return blank email",
		data:     AddAdmin{"John Name", " ", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	invalidEmail1 := SampleAdmin{
		name:     "should return invalid email",
		data:     AddAdmin{"John Doe", "user@gmailcom", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	invalidEmail2 := SampleAdmin{
		name:     "should return invalid email",
		data:     AddAdmin{"John Name", "usergmail.com", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	shortEmail := SampleAdmin{
		name:     "should return too short email",
		data:     AddAdmin{"John Name", "u@gl.c", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	longEmail := SampleAdmin{
		name:     "should return too long email",
		data:     AddAdmin{"John Name", "uasdjdsksksksksksk@gmaillll.coooom", "0712345689", "Acx@1234", "this is image", "role"},
		expected: false,
	}

	// cases for phone number
	tooShortCell := SampleAdmin{
		name:     "should return too short cell phone number",
		data:     AddAdmin{"John Name", "user@gmail.com", "071234689", "Acx@1234", "this is image", "role"},
		expected: false,
	}

	tooLongCell := SampleAdmin{
		name:     "should return too long cell phone number",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345285858689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	letterInCell := SampleAdmin{
		name:     "should return a letter present phone number",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345ds689", "Acx@1234", "this is image", "role"},
		expected: false,
	}
	specialCharInCell := SampleAdmin{
		name:     "should return special char in phone number",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345@#689", "Acx@1234", "this is image", "role"},
		expected: false,
	}

	// casses for password
	noCapsPassword := SampleAdmin{
		name:     "should return no caps in",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345689", "acx@1234", "this is image", "role"},
		expected: false,
	}
	noCharPassword := SampleAdmin{
		name:     "should return no special character in password",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345689", "Acxx1234", "this is image", "role"},
		expected: false,
	}
	noNumPassword := SampleAdmin{
		name:     "should return too long cell phone number",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345689", "AcxxxxcA", "this is image", "role"},
		expected: false,
	}
	tooShortPassword := SampleAdmin{
		name:     "should return too long cell phone number",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345689", "A@1x", "this is image", "role"},
		expected: false,
	}
	tooShortRole := SampleAdmin{
		name:     "should return too long cell phone number",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345689", "A@1x", "this is image", "ro"},
		expected: false,
	}
	cases := []SampleAdmin{
		emptyName,
		shortName,
		blankName,
		longName,
		numInName,
		specialCharInName,
		emptyEmail,
		blankEmail,
		invalidEmail1,
		invalidEmail2,
		shortEmail,
		longEmail,
		tooShortCell,
		tooLongCell,
		letterInCell,
		specialCharInCell,
		noCapsPassword,
		noCharPassword,
		noNumPassword,
		tooShortPassword,
		tooShortRole,
	}

	for _, item := range cases {
		result, _ := ValidateRegisterInput(&item.data)

		if item.expected != result {
			t.Errorf("error occurred!! test %s expected %v but found %v", item.name, item.expected, result)
		}

	}
}
