package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
	"project1.com/project/logsetup"
)

type Logininfo struct {
	Username string
	Password string
}

//checks user is present and allows user to login is password matches in db
func Login(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	file, flag := logsetup.Logfile(w, res)
	defer file.Close()
	if flag {
		return
	}
	log.SetOutput(file) //setting output destination
	detail, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		error(res, w)
		log.Print(err)
		return
	}
	var det Logininfo
	var det1 Registration
	err = json.Unmarshal(detail, &det) //convert json to struct
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		error(res, w)
		log.Print(err)
		return
	}
	err = db.Model(&det1).Where("username=?", det.Username).Select()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res["message"] = "No User Found"
		error(res, w)
		log.Print(err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(det1.Password), []byte(det.Password)) //decrypt password
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized) //status code for unathorization
		res["message"] = "Entered password is wrong!!!"
		error(res, w)
		log.Print(err)
		return
	} else {
		str := fmt.Sprintf("%s Welcome", det.Username)
		res["message"] = str
		error(res, w)
		return
	}

}
