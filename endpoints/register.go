package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strings"

	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
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

//registers the user data to the table registration in database
func PostRegistration(w http.ResponseWriter, r *http.Request, db *pg.DB) {

	detail, err := io.ReadAll(r.Body) //reads the request body and returns byte value
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		jsonstr, _ := json.Marshal(res)
		w.Write(jsonstr)
		return
	}
	var det Registration
	err = json.Unmarshal(detail, &det) //converts json format to struct
	if err != nil {
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		jsonstr, _ := json.Marshal(res)
		w.Write(jsonstr)
		return
	}
	if det.Firstname == "" || det.Lastname == "" || det.Username == "" || det.Password == "" || det.Email == "" {
		res["message"] = "Please Enter all the details"
		jsonstr, _ := json.Marshal(res)
		w.Write(jsonstr)
		return
	}

	_, err = mail.ParseAddress(det.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter valid email ID"
		jsonstr, _ := json.Marshal(res)
		w.Write(jsonstr)
		return
	}

	/*passwordvalidator := validator.New(validator.MinLength(8, errors.New("too short")),
		validator.ContainsAtLeast("abcdefghijklmnopqrstuvwxyz", 2, errors.New("Contain atleast 2 lowercase")),
		validator.ContainsAtLeast("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 2, errors.New("Contain atleast 2 Uppercase")),
		validator.ContainsAtLeast("1234567890", 1, errors.New("Contain atleast 1 Numbers")),
		validator.ContainsAtLeast("@$!%*#?&", 1, errors.New("Should Contain atleast 1 Special Charecter")))
	err = passwordvalidator.Validate(det.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		str := fmt.Sprintf("%s,Note:Password Should contain Atleast 2 Uppercase,Lowercase And 1 Numbers,Special Char", err.Error())
		res["message"] = str
		jsonstr, _ := json.Marshal(res)
		w.Write(jsonstr)
		return
	}*/
	//re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	/*_, err = regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", det.Email)
	if err != nil {
		error(err)
		return
	}*/
	//re := regexp.MustCompile("^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9])(?=.*[@$!%*#?&])[A-Za-z0-9@$!%*#?&]{8,}$")
	/*if _, err = regexp.MatchString("^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[*.!@$%^&(){}[]:;<>,.?/~_+-=|\\]).{8,32}$", det.Password); err != nil {
		//error(err)
		fmt.Println(err)
		return
	}*/

	bytes, err := bcrypt.GenerateFromPassword([]byte(det.Password), 10) //Encryption of password feild to bytes
	if err != nil {
		res["message"] = "Something wrong in backend...Failed to encrypt password"
		jsonstr, _ := json.Marshal(res)
		w.Write(jsonstr)
		return
	}
	det.Password = string(bytes)     //converts byte to string and update the feild of password
	_, err = db.Model(&det).Insert() //query to Insert the data into database

	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		str := err.Error()
		last := str[strings.LastIndex(str, " ")+2 : strings.LastIndex(str, " ")+25]
		if last == "registrations_email_key" {
			res["message"] = "Email-Id is already registered"
			jsonstr, _ := json.Marshal(res)
			w.Write(jsonstr)
		} else {
			res["message"] = "Username is already registered"
			jsonstr, _ := json.Marshal(res)
			w.Write(jsonstr)
		}

	} else {
		w.WriteHeader(http.StatusCreated)
		str := fmt.Sprintf("%s successfully registered", det.Username)
		res["message"] = str
		jsonstr, _ := json.Marshal(res)
		//fmt.Fprint(w, string(jsonstr))
		w.Write(jsonstr)

	}
}
