package validation

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/go-passwd/validator"
	valid "github.com/go-validator/validator"
	log "github.com/sirupsen/logrus"
	"project1.com/project/utility"
)

/*func display(res map[string]string, w http.ResponseWriter) {
	jsonstr, _ := json.Marshal(res)
	w.Write(jsonstr)
}*/
//validation of feilds wheather empty or not
func FeildValidation(det interface{}, w http.ResponseWriter, res map[string]string) error {
	if err := valid.Validate(det); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter all the details"
		utility.Display(res, w)
		log.Warn(res["message"])
		return err
	}
	return nil
}

//validation of password
func Passwordvalidation(res map[string]string, password string, w http.ResponseWriter) bool {
	passwordvalidator := validator.New(validator.MinLength(8, errors.New("too short")),
		validator.ContainsAtLeast("abcdefghijklmnopqrstuvwxyz", 2, errors.New("Contain atleast 2 lowercase")),
		validator.ContainsAtLeast("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 2, errors.New("Contain atleast 2 Uppercase")),
		validator.ContainsAtLeast("1234567890", 1, errors.New("Contain atleast 1 Numbers")),
		validator.ContainsAtLeast("@$!%*#?&", 1, errors.New("Should Contain atleast 1 Special Charecter")))
	err := passwordvalidator.Validate(password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		str := fmt.Sprintf("%s,Note:Password Should contain Atleast 2 Uppercase,Lowercase And 1 Number,Special Char", err.Error())
		res["message"] = str
		utility.Display(res, w)
		log.Warn(err)
		return true
	}
	return false
}

//Email valdation
func Emailvalidation(res map[string]string, email string, w http.ResponseWriter) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter valid email ID"
		utility.Display(res, w)
		log.Warn(err)
		return true
	}
	return false
}
