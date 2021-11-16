package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
	"project1.com/project/passemailvalidation"
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

//displays errors to user end
func error(res map[string]string, w http.ResponseWriter) {
	jsonstr, _ := json.Marshal(res)
	w.Write(jsonstr)
}

func logfile(w http.ResponseWriter) (*os.File, bool) {
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) //opening a log file
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant Open Log file"
		error(res, w)
		return nil, true
	}
	return file, false
}

//registers the user data to the table registration in database
func PostRegistration(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	file, flag := logfile(w)
	if flag {
		return
	}
	log.SetOutput(file)               //setting output destination
	detail, err := io.ReadAll(r.Body) //reads the request body and returns byte value
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		error(res, w)
		log.Print(err)
		return
	}
	var det Registration
	err = json.Unmarshal(detail, &det) //converts json format to struct
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		error(res, w)
		log.Print(err)
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
	if flag := passemailvalidation.Passwordvalidation(res, det.Password, w); flag {
		fmt.Print("hello")
		return
	}
	if flag := passemailvalidation.Emailvalidation(res, det.Email, w); flag {
		return
	}
	//encrption of password
	bytes, err := bcrypt.GenerateFromPassword([]byte(det.Password), 10) //Encryption of password feild to bytes
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend...Failed to encrypt password"
		error(res, w)
		log.Print(err)
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
		log.Print(err)
	} else {
		w.WriteHeader(http.StatusCreated)
		str := fmt.Sprintf("%s successfully registered", det.Username)
		res["message"] = str
		error(res, w)
	}
}
