package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-pg/pg"
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
}

var res = map[string]string{"message": ""}
var bytes []byte

//displays errors to user end
func error(res map[string]string, w http.ResponseWriter) {
	jsonstr, _ := json.Marshal(res)
	w.Write(jsonstr)
}

//registers the user data to the table registration in database
func PostRegistration(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	file, flag := logsetup.Logfile(w, res)
	defer file.Close()
	if flag {
		return
	}
	log.SetOutput(file)               //setting output destination
	detail, err := io.ReadAll(r.Body) //reads the request body and returns byte value
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		error(res, w)
		log.Println(err)
		return
	}
	var det Registration
	err = json.Unmarshal(detail, &det) //converts json format to struct
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		error(res, w)
		log.Println(err)
		return
	}
	//validation
	if det.Firstname == "" || det.Firstname == " " || det.Lastname == " " || det.Lastname == "" || det.Username == " " ||
		det.Username == "" || det.Password == " " || det.Password == "" || det.Email == " " || det.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter all the details"
		error(res, w)
		return
	}
	if flag := validation.Passwordvalidation(res, det.Password, w); flag {
		return
	}
	if flag := validation.Emailvalidation(res, det.Email, w); flag {
		return
	}
	//encrption of password
	if bytes = validation.Encrption(det.Password, w, res); bytes == nil {
		return
	}
	det.Password = string(bytes)     //converts byte to string and update the feild of password
	_, err = db.Model(&det).Insert() //query to Insert the data into database
	//checks wheather username or email exists
	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		str := err.Error()
		last := str[strings.LastIndex(str, " ")+2 : strings.LastIndex(str, " ")+25]
		if last == "registrations_email_key" {
			res["message"] = "Email-Id is already registered"
			error(res, w)
		} else {
			res["message"] = "Username is already registered"
			error(res, w)
		}
		log.Println(err)
	} else {
		w.WriteHeader(http.StatusCreated)
		str := fmt.Sprintf("%s successfully registered", det.Username)
		res["message"] = str
		error(res, w)
	}
}
