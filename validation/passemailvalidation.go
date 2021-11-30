package validation

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/go-passwd/validator"
	valid "github.com/go-validator/validator"
)

//validation of feilds wheather empty or not
func FeildValidation(det interface{}) error {
	if err := valid.Validate(det); err != nil {
		return err
	}
	return nil
}

//validation of password
func PasswordValidation(password string) error {
	passwordvalidator := validator.New(validator.MinLength(8, errors.New("too short")),
		validator.ContainsAtLeast("abcdefghijklmnopqrstuvwxyz", 2, errors.New("Contain atleast 2 lowercase")),
		validator.ContainsAtLeast("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 2, errors.New("Contain atleast 2 Uppercase")),
		validator.ContainsAtLeast("1234567890", 1, errors.New("Contain atleast 1 Numbers")),
		validator.ContainsAtLeast("@$!%*#?&", 1, errors.New("Should Contain atleast 1 Special Charecter")))
	if err := passwordvalidator.Validate(password); err != nil {
		return err
	}
	return nil
}

//Email valdation
func EmailValidation(res map[string]string, email string, w http.ResponseWriter) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}
	return nil
}
