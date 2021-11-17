package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"

	"github.com/go-passwd/validator"
	"project1.com/project/logsetup"
)

func error(res map[string]string, w http.ResponseWriter) {
	jsonstr, _ := json.Marshal(res)
	w.Write(jsonstr)
}
func Passwordvalidation(res map[string]string, password string, w http.ResponseWriter) bool {
	file, flag := logsetup.Logfile(w, res)
	defer file.Close()
	if flag {
		return true
	}
	log.SetOutput(file)
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
		error(res, w)
		log.Println(err.Error())
		return true
	}
	return false
}
func Emailvalidation(res map[string]string, email string, w http.ResponseWriter) bool {
	file, flag := logsetup.Logfile(w, res)
	defer file.Close()
	if flag {
		return true
	}
	log.SetOutput(file)
	//email validation
	_, err := mail.ParseAddress(email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter valid email ID"
		error(res, w)
		log.Print(err)
		return true
	}
	return false
}
