package endpoints

import (
	"fmt"

	"net/http"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	read "project1.com/project/endpoints/readrequestbody"
	"project1.com/project/logsetup"
	"project1.com/project/utility"
	"project1.com/project/validation"
)

type Logininfo struct {
	Username string `validate:"nonzero"`
	Password string `validate:"nonzero"`
}

//checks user is present and allows user to login is password matches in db
func Login(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	file, flag := logsetup.Logfile(w, res)
	if flag {
		return
	}
	defer file.Close()
	log.SetOutput(file) //setting output destination
	var det *Logininfo
	var det1 Registration
	if err := read.Readbody(r, w, res, &det); err != nil {
		return
	}
	//validation
	if err := validation.FeildValidation(det); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter all the details"
		utility.Display(res, w)
		log.Warn(res["message"])
		return
	}
	err := db.Model(&det1).Where("username=?", det.Username).Select()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res["message"] = "No User Found"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(det1.Password), []byte(det.Password)) //decrypt password
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized) //status code for unathorization
		res["message"] = "Entered password is wrong!!!"
		utility.Display(res, w)
		log.Error(err)
	} else {
		w.WriteHeader(http.StatusFound)
		str := fmt.Sprintf("%s Welcome", det.Username)
		res["message"] = str
		utility.Display(res, w)
		log.Info(str)
	}

}
