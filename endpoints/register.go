package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

//diplays the error to the user end
func badrequesterror(err interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Println(err)
}
func error(err interface{}) {
	fmt.Println(err)
}

//registers the user data to the table registration in database
func PostRegistration(w http.ResponseWriter, r *http.Request, db *pg.DB) {

	detail, err := io.ReadAll(r.Body) //reads the request body and returns byte value
	if err != nil {
		badrequesterror(err, w)
		return
	}
	var det Registration
	err = json.Unmarshal(detail, &det) //converts json format to struct
	if err != nil {
		error(err)
		return
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(det.Password), 10) //Encryption of password feild to bytes
	if err != nil {
		error(err)
		return
	}
	det.Password = string(bytes)     //converts byte to string and update the feild of password
	_, err = db.Model(&det).Insert() //query to Insert the data into database

	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)

	} else {
		w.WriteHeader(http.StatusCreated)
		str := fmt.Sprintf("%s successfully registered", det.Username)
		res := map[string]string{"message": str}
		jsonstr, err := json.Marshal(res)
		if err != nil {
			error(err)
			return
		}
		//fmt.Fprint(w, string(jsonstr))
		w.Write(jsonstr)
	}
}
