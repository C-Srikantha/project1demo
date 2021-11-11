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
func error(err interface{}) {
	if err != nil {
		fmt.Println(err)
	}
}

//registers the user data to the table registration in database
func PostRegistration(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	detail, err := io.ReadAll(r.Body) //reads the request body and returns byte value
	error(err)
	var det Registration
	err = json.Unmarshal(detail, &det) //converts json format to struct
	error(err)
	bytes, err := bcrypt.GenerateFromPassword([]byte(det.Password), 10) //Encryption of password feild to bytes
	error(err)
	det.Password = string(bytes)     //converts byte to string and update the feild of password
	_, err = db.Model(&det).Insert() //query to Insert the data into database
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Fprintf(w, "%s succesfully registerd", det.Username) //Prints successful message to user end
	}

}
