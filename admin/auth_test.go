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

}
