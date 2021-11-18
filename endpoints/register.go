package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"project1.com/project/logsetup"
	"project1.com/project/validation"
)

type Registration struct {
	Id        int
	Firstname string
	Lastname  string
	Username  string
	Password  string
	Email     string
	Otp       string
}

var res = map[string]string{"message": ""}
var bytepass []byte

//displays errors to user end
func display(res map[string]string, w http.ResponseWriter) {
	jsonstr, _ := json.Marshal(res)
	w.Write(jsonstr)
}

//registers the user data to the table registration in database
func PostRegistration(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	file, flag := logsetup.Logfile(w, res)
	if flag {
		return
	}
	//defer file.Close()
	log.SetOutput(file)               //setting output destination
	detail, err := io.ReadAll(r.Body) //reads the request body and returns byte value
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		display(res, w)
		log.Error(err)
		return
	}
	var det Registration
	err = json.Unmarshal(detail, &det) //converts json format to struct
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		display(res, w)
		log.Error(err)
		return
	}
	//validation
	if det.Firstname == "" || det.Firstname == " " || det.Lastname == " " || det.Lastname == "" || det.Username == " " ||
		det.Username == "" || det.Password == " " || det.Password == "" || det.Email == " " || det.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter all the details"
		display(res, w)
		log.Warn(res["message"])
		return
	}
	if flag := validation.Passwordvalidation(res, det.Password, w); flag {
		return
	}
	if flag := validation.Emailvalidation(res, det.Email, w); flag {
		return
	}
	//encrption of password
	if bytepass = validation.Encrption(det.Password, w, res); bytepass == nil {
		return
	}
	det.Password = string(bytepass)  //converts byte to string and update the feild of password
	_, err = db.Model(&det).Insert() //query to Insert the data into database
	//checks wheather username or email exists
	if err != nil {
		str := err.Error()
		w.WriteHeader(http.StatusAlreadyReported)
		last := str[strings.LastIndex(str, " ")+2 : strings.LastIndex(str, " ")+25]
		if last == "registrations_email_key" {
			res["message"] = "Email-Id is already registered"
			display(res, w)
			log.Warn(res["message"])
		} else {
			res["message"] = "Username is already registered"
			display(res, w)
			log.Warn(res["message"])
		}
		log.Error(err)
	} else {
		w.WriteHeader(http.StatusCreated)
		str := fmt.Sprintf("%s successfully registered", det.Username)
		res["message"] = str
		display(res, w)
		log.Info(str)
	}
}
