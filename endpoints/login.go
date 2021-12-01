package endpoints

import (
	"fmt"
	"os"

	"net/http"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"project1.com/project/createtable"
	read "project1.com/project/endpoints/readrequestbody"
	"project1.com/project/utility"
	"project1.com/project/validation"
)

type Logininfo struct {
	Username string `validate:"nonzero"`
	Password string `validate:"nonzero"`
}

//Login function checks user is present and allows user to login is password matches in db
func Login(w http.ResponseWriter, r *http.Request, db *pg.DB, file *os.File) {

	log.SetOutput(file) //setting output destination
	var det *Logininfo
	var det1 createtable.Registration
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
	//validation of feilds is empty or not
	if err := validation.FeildValidation(det); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter all the details"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	//Username exist or not
	err = db.Model(&det1).Where("username=?", det.Username).Select()
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
