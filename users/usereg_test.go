package users

import (
	"testing"
)

// testt for register data
func TestUserRegData(t *testing.T) {
	// cases for empty values
	invalidFnameLength := RegisterInput{"ia", "Does", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "0719722292", "Pass@1234"}
	invalidLnameLength := RegisterInput{"John", "Does", "do", "user image", "Nyeri", "JohnDoe@gmail.com", "0719722292", "Pass@1234"}
	invalidLocationLength := RegisterInput{"John", "Does", "Doe", "user image", "", "JohnDoe@gmail.com", "0719722292", "Pass@1234"}
	invalidEmailLength := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "", "0719722292", "Pass@1234"}
	invalidPasswordLength := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "07122292", ""}
	invalidPhoneNumberLength := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "", "Pass@1234"}

	// cases for empty values
	emptyFname := RegisterInput{"   ", "Does", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "0719722292", "Pass@1234"}
	emptyLname := RegisterInput{"John", "Does", "   ", "user image", "Nyeri", "JohnDoe@gmail.com", "0719722292", "Pass@1234"}
	emptyEmail := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "  ", "0719722292", "Pass@1234"}
	emptyLocation := RegisterInput{"John", "Does", "Doe", "user image", "  ", "JohnDoe@gmail.com", "0719722292", "Pass@1234"}
	emptyPhoneNumber := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "   ", "Pass@1234"}
	emptyPassword := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "0712345689", "  "}

	// special cases for password
	noCapsInPassword := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "0712345689", "pass@1234"}
	noNumInPassword := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "0712345689", "Pass@Pass"}
	noSpecialCharInPassword := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "0712345689", "Pass1234"}

	// cases for invalid characters
	invalidEmailType := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "johngmail.com", "0719722292", "Pass@1234"}
	invalidEmailType2 := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "john@gmailcom", "0719722292", "Pass@1234"}
	invalidCharactersInEmail := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "jo<>#&*%(hn@gmail.com", "0719722292", "Pass@1234"}
	invalidCharactersInFname := RegisterInput{"@$$%&*)(*)*_*", "Does", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "07122292", "Pass@1234"}
	invalidCharactersInMname := RegisterInput{"@$$%&*)(*)*_*", "@$$%&*)(*)*_*", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "07122292", "Pass@1234"}
	invalidCharactersInLname := RegisterInput{"John", "Does", "@john$$%&*)(*)*_*", "user image", "Nyeri", "JohnDoe@gmail.com", "07122292", "Pass@1234"}
	invalidCharactersInPhoneNumber := RegisterInput{"John", "Does", "@john$$%&*)(*)*_*", "user image", "Nyeri", "JohnDoe@gmail.com", "0712345#%689", "Pass@1234"}
	lettersInPhoneNumber := RegisterInput{"John", "Does", "@john$$%&*)(*)*_*", "user image", "Nyeri", "JohnDoe@gmail.com", "071234RAD689", "Pass@1234"}
	invalidCharactersInLocation := RegisterInput{"John", "Does", "@john$$*)(*)*_*", "user image", "Ny$#*^ri", "JohnDoe@gmail.com", "071234RAD689", "Pass@1234"}

	validData := RegisterInput{"John", "Does", "Doe", "user image", "Nyeri", "JohnDoe@gmail.com", "0719722292", "Pass@1234"}

	// table driven
	cases := []struct {
		name string
		user RegisterInput
		want bool
	}{
		// test cases for empty values
		{
			name: "should return invalid length for first name",
			user: invalidFnameLength,
			want: false,
		},
		{
			name: "should return invalid length for last name",
			user: invalidLnameLength,
			want: false,
		},
		{
			name: "should return invalid length for location",
			user: invalidLocationLength,
			want: false,
		},
		{
			name: "should return invalid length email",
			user: invalidEmailLength,
			want: false,
		},
		{
			name: "should return invalid length for password",
			user: invalidPasswordLength,
			want: false,
		},
		{
			name: "should return invalid length for phone number",
			user: invalidPhoneNumberLength,
			want: false,
		},
		// cases for empty values
		{
			name: "should return empty first name",
			user: emptyFname,
			want: false,
		},
		{
			name: "should return empty last name",
			user: emptyLname,
			want: false,
		},
		{
			name: "should return empty location",
			user: emptyLocation,
			want: false,
		},
		{
			name: "should return empty email",
			user: emptyEmail,
			want: false,
		},
		{
			name: "should return empty phone number",
			user: emptyPhoneNumber,
			want: false,
		},
		{
			name: "should return empty password",
			user: emptyPassword,
			want: false,
		},
		// cases invalid characters
		{
			name: "should return invalid email type",
			user: invalidEmailType,
			want: false,
		},
		{
			name: "should return invalid email type .",
			user: invalidEmailType2,
			want: false,
		},
		{
			name: "should return invalid characters in email",
			user: invalidCharactersInEmail,
			want: false,
		},
		{
			name: "should return invalid characters in first name",
			user: invalidCharactersInFname,
			want: false,
		},
		{
			name: "should return invalid characters in middle name",
			user: invalidCharactersInMname,
			want: false,
		},
		{
			name: "should return invalid characters in last name",
			user: invalidCharactersInLname,
			want: false,
		},
		{
			name: "should return invalid characters in phone number",
			user: invalidCharactersInPhoneNumber,
			want: false,
		},
		{
			name: "should return no characters in phone number",
			user: lettersInPhoneNumber,
			want: false,
		},
		{
			name: "should return no caps in password",
			user: noCapsInPassword,
			want: false,
		},
		// special password cases
		{
			name: "should return no number in password",
			user: noNumInPassword,
			want: false,
		},
		{
			name: "should return no special characters in password",
			user: noSpecialCharInPassword,
			want: false,
		},
		{
			name: "should return invalid characters in location",
			user: invalidCharactersInLocation,
			want: false,
		},
		{
			name: "should return true and valid data",
			user: validData,
			want: true,
		},
	}

	for _, input := range cases {
		result, err := ValidateRegisterInput(&input.user)
		if result != input.want {
			t.Errorf("test failed %s", err)
		}
	}
	t.Logf("test is a success")
}

// test for login data

func TestLoginData(t *testing.T) {
	// email test
	invalidEmailLength := LoginInput{" ", "Pass@1234"}
	invalidEmailLength2 := LoginInput{"", "Pass@1234"}
	invalidEmailType := LoginInput{"john@gmailcom", "Pass@1234"}
	invalidEmailType2 := LoginInput{"johngmail.com", "Pass@1234"}
	invalidSpecialCharInEmail := LoginInput{"johng@$%*&^mail.com", "Pass@1234"}

	// password tests
	invalidPasswordLength := LoginInput{"john@gmailcom", " "}
	invalidPasswordLength2 := LoginInput{"john@gmailcom", " "}
	noCaps := LoginInput{"john@gmailcom", "pass@1234"}
	noSpecialChar := LoginInput{"johngmail.com", "Pass1234"}

	// valid data
	validData := LoginInput{"john@gmail.com", "Pass1234"}

	cases := []struct {
		name string
		data LoginInput
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
