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
		name:     "admin name is empty",
		data:     AddAdmin{"", "user@gmail.com", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	blankName := SampleAdmin{
		name:     "should return blank name",
		data:     AddAdmin{"   ", "user@gmail.com", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	shortName := SampleAdmin{
		name:     "should return too short name",
		data:     AddAdmin{"as", "user@gmail.com", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	longName := SampleAdmin{
		name:     "should return too long name",
		data:     AddAdmin{"this is a very long name that is not allowed", "user@gmail.com", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	numInName := SampleAdmin{
		name:     "should return number are not allowed",
		data:     AddAdmin{"this1235", "user@gmail.com", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	specialCharInName := SampleAdmin{
		name:     "should return name should not contain special characters",
		data:     AddAdmin{"shoulname@#", "user@gmail.com", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}

	// cases for email
	emptyEmail := SampleAdmin{
		name:     "should return empty email",
		data:     AddAdmin{"this is a very long name", "", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	blankEmail := SampleAdmin{
		name:     "should return blank email",
		data:     AddAdmin{"John Name", " ", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	invalidEmail1 := SampleAdmin{
		name:     "should return invalid email",
		data:     AddAdmin{"John Doe", "user@gmailcom", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	invalidEmail2 := SampleAdmin{
		name:     "should return invalid email",
		data:     AddAdmin{"John Name", "usergmail.com", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	shortEmail := SampleAdmin{
		name:     "should return too short email",
		data:     AddAdmin{"John Name", "u@gl.c", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	longEmail := SampleAdmin{
		name:     "should return too long email",
		data:     AddAdmin{"John Name", "uasdjdsksksksksksk@gmaillll.coooom", "0712345689", "Pass@1234", "this is image", "role"},
		expected: false,
	}

	// cases for phone number
	tooShortCell := SampleAdmin{
		name:     "should return too short cell phone number",
		data:     AddAdmin{"John Name", "user@gmail.com", "071234689", "Pass@1234", "this is image", "role"},
		expected: false,
	}

	tooLongCell := SampleAdmin{
		name:     "should return too long cell phone number",
		data:     AddAdmin{"John Name", "user@gmail.com", "071234568912345698", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	letterInCell := SampleAdmin{
		name:     "should return a letter present phone number",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345ds689", "Pass@1234", "this is image", "role"},
		expected: false,
	}
	specialCharInCell := SampleAdmin{
		name:     "should return special char in phone number",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345@#689", "Pass@1234", "this is image", "role"},
		expected: false,
	}

	// casses for password
	noCapsPassword := SampleAdmin{
		name:     "should return no caps in",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345689", "pass@1234", "this is image", "role"},
		expected: false,
	}
	noCharPassword := SampleAdmin{
		name:     "passsword has no special characters",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345689", "Pass1234", "this is image", "role"},
		expected: false,
	}
	noNumPassword := SampleAdmin{
		name:     "password has no numbers",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345689", "Pass@Pass", "this is image", "role"},
		expected: false,
	}
	tooShortPassword := SampleAdmin{
		name:     "password too short",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345689", "A@1x", "this is image", "role"},
		expected: false,
	}
	tooShortRole := SampleAdmin{
		name:     "role too short",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345689", "Pass@1234", "this is image", "ro"},
		expected: false,
	}

	validData := SampleAdmin{
		name:     "should return valid data",
		data:     AddAdmin{"John Name", "user@gmail.com", "0712345689", "Pass@1234", "this is image", "role"},
		expected: true,
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
		validData,
	}

	for _, item := range cases {
		result, _ := ValidateRegisterInput(&item.data)

		if item.expected != result {
			t.Errorf("error occurred!! test %s expected %v but found %v", item.name, item.expected, result)
		}

	}
}

func TestLoginData(t *testing.T) {
	// email test
	invalidEmailLength := AdminLogin{" ", "Pass@1234"}
	invalidEmailLength2 := AdminLogin{"", "Pass@1234"}
	invalidEmailType := AdminLogin{"john@gmailcom", "Pass@1234"}
	invalidEmailType2 := AdminLogin{"johngmail.com", "Pass@1234"}
	invalidSpecialCharInEmail := AdminLogin{"johng@$%*&^mail.com", "Pass@1234"}

	// password tests
	invalidPasswordLength := AdminLogin{"john@gmailcom", " "}
	invalidPasswordLength2 := AdminLogin{"john@gmailcom", " "}
	noCaps := AdminLogin{"john@gmailcom", "pass@1234"}
	noSpecialChar := AdminLogin{"johngmail.com", "Pass1234"}

	// valid data
	validData := AdminLogin{"john@gmail.com", "Pass1234"}

	cases := []struct {
		name string
		data AdminLogin
		want bool
	}{
		{
			name: "should return invalid length email",
			data: invalidEmailLength,
			want: false,
		},
		{
			name: "should return invalid length email",
			data: invalidEmailLength2,
			want: false,
		},
		{
			name: "should return invalid type of email",
			data: invalidEmailType,
			want: false,
		},
		{
			name: "should return invalid type of email",
			data: invalidEmailType2,
			want: false,
		},
		{
			name: "should return invalid characters in email",
			data: invalidSpecialCharInEmail,
			want: false,
		},

		// cases for password
		{
			name: "should return invalid length for password",
			data: invalidPasswordLength,
			want: false,
		},
		{
			name: "should return invalid length for password",
			data: invalidPasswordLength2,
			want: false,
		},
		{
			name: "should return no caps in password",
			data: noCaps,
			want: false,
		},
		{
			name: "should return no special characters in password",
			data: noSpecialChar,
			want: false,
		},
		{
			name: "test passed",
			data: validData,
			want: false,
		},
	}

	for _, input := range cases {
		result, err := ValidateLoginInput(&input.data)
		if result != input.want {
			t.Errorf("test failed: %s%v", err, input.data)
		}
	}
	t.Logf("test is successful")
}
