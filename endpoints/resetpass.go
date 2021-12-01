package endpoints

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"project1.com/project/createtable"
	read "project1.com/project/endpoints/readrequestbody"
	"project1.com/project/utility"
	"project1.com/project/validation"
)

type Resetpass struct {
	Username    string `validate:"nonzero"`
	Otp         string `validate:"nonzero"`
	Newpassword string `validate:"nonzero"`
}

func Reset(w http.ResponseWriter, r *http.Request, db *pg.DB, file *os.File) {
	log.SetOutput(file) //setting output destination
	var det *Resetpass
	var det1 createtable.Registration
	//reading body and store values to var det
	if err := read.Readbody(r, w, res, &det); err != nil {
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
	if err := validation.PasswordValidation(det.Newpassword); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		str := fmt.Sprintf("%s,Note:Password Should contain Atleast 2 Uppercase,Lowercase And 1 Number,Special Char", err.Error())
		res["message"] = str
		utility.Display(res, w)
		log.Error(err)
		return
	}
	//Username exist or not
	err := db.Model(&det1).Where("username=?", det.Username).Select() //query to check username exist or not
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please enter valid Username "
		utility.Display(res, w)
		log.Error(err)
		return
	}
	//checks otp matches with database
	err = bcrypt.CompareHashAndPassword([]byte(det1.Otp), []byte(det.Otp))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized) //status code for unathorization
		res["message"] = "OTP Entered is wrong!!!"
		utility.Display(res, w)
		log.Error(err)
		return
	}

	//Encryption of password
	if bytepass, err = utility.Encrption(det.Newpassword); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend...Failed to encrypt password"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	//updating password into database
	_, err = db.Model(&det1).Set("password=?", string(bytepass)).Where("username=?", det.Username).Update()
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		res["message"] = "Password reset failed"
		utility.Display(res, w)
		log.Error(err)
	} else {
		w.WriteHeader(http.StatusCreated)
		res["message"] = "Password reset success"
		utility.Display(res, w)
		log.Info(res["message"])
	}

}
