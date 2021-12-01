package endpoints

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"

	read "project1.com/project/endpoints/readrequestbody"
	"project1.com/project/utility"
	"project1.com/project/validation"
)

type Registration struct {
	Id        int
	Firstname string `validate:"nonzero"`
	Lastname  string `validate:"nonzero"`
	Username  string `validate:"nonzero"`
	Password  string `validate:"nonzero"`
	Email     string `validate:"nonzero"`
	Otp       string
}

var res = map[string]string{"message": ""}
var bytepass []byte

//registers the user data to the table registration in database
func PostRegistration(w http.ResponseWriter, r *http.Request, db *pg.DB, file *os.File) {
	/*	file, flag := logsetup.Logfile(w, res)
		if flag {
			return
		}
	defer file.Close()*/
	log.SetOutput(file) //setting log output destination
	var det *Registration
	if err := read.Readbody(r, w, res, &det); err != nil { //calling Readbody for reading requestbody
		return
	}
	if err := validation.FeildValidation(det); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter all the details"
		utility.Display(res, w)
		log.Warn(res["message"])
		return
	}
	if err := validation.PasswordValidation(det.Password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		str := fmt.Sprintf("%s,Note:Password Should contain Atleast 2 Uppercase,Lowercase And 1 Number,Special Char", err.Error())
		res["message"] = str
		utility.Display(res, w)
		log.Error(err)
		return
	}
	if err := validation.EmailValidation(res, det.Email, w); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter valid email ID"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	//encrption of password
	if bytepass, err := utility.Encrption(det.Password); bytepass == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend...Failed to encrypt password"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	det.Password = string(bytepass)  //converts byte to string and update the feild of password
	_, err := db.Model(det).Insert() //query to Insert the data into database
	//checks wheather username or email exists
	if err != nil {
		str := err.Error()
		w.WriteHeader(http.StatusAlreadyReported)
		last := str[strings.LastIndex(str, " ")+2 : strings.LastIndex(str, " ")+25]
		if last == "registrations_email_key" {
			res["message"] = "Email-Id is already registered"
			utility.Display(res, w)
			log.Warn(res["message"])
		} else {
			res["message"] = "Username is already registered"
			utility.Display(res, w)
			log.Warn(res["message"])
		}
		log.Error(err)
	} else {
		w.WriteHeader(http.StatusCreated)
		str := fmt.Sprintf("%s successfully registered", det.Username)
		res["message"] = str
		utility.Display(res, w)
		log.Info(str)
	}
}
