package endpoints

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"

	read "project1.com/project/endpoints/readrequestbody"
	"project1.com/project/enti"
	"project1.com/project/utility"
	"project1.com/project/validation"
)

var res = map[string]string{"message": ""}
var bytepass []byte

//PostRegistration function registers the user data to the table registration in database
func PostRegistration(w http.ResponseWriter, r *http.Request, db *pg.DB, file *os.File) {

	log.SetOutput(file) //setting log output destination
	var det *enti.Registration
	detail, err := read.ReadBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	if err := read.Convert(detail, &det); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	//validation of feilds if empty or not
	if err := validation.FeildValidation(det); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter all the details"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	//validation of password
	if err := validation.PasswordValidation(det.Password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		str := fmt.Sprintf("%s,Note:Password Should contain Atleast 2 Uppercase,Lowercase And 1 Number,Special Char", err.Error())
		res["message"] = str
		utility.Display(res, w)
		log.Error(err)
		return
	}
	//validation of Email
	if err := validation.EmailValidation(res, det.Email, w); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter valid email ID"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	//encrption of password
	bytepass, err := utility.Encrption(det.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend...Failed to encrypt password"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	det.Password = string(bytepass) //converts byte to string and update the feild of password
	_, err = db.Model(det).Insert() //query to Insert the data into database
	//checks wheather username or email exists
	if err != nil {
		str := err.Error()
		w.WriteHeader(http.StatusAlreadyReported)
		last := str[strings.LastIndex(str, " ")+2 : strings.LastIndex(str, " ")+25]
		if last == "registrations_email_key" {
			res["message"] = "Email-Id is already registered"
			utility.Display(res, w)
		} else {
			res["message"] = "Username is already registered"
			utility.Display(res, w)
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
