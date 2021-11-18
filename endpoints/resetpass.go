package endpoints

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"project1.com/project/logsetup"
	"project1.com/project/validation"
)

type Resetpass struct {
	Username    string
	Otp         string
	Newpassword string
}

func Reset(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	file, flag := logsetup.Logfile(w, res)
	if flag {
		return
	}
	defer file.Close()
	log.SetOutput(file) //setting output destination
	//reading body
	detail, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		display(res, w)
		log.Error(err)
		return
	}
	var det Resetpass
	var det1 Registration
	err = json.Unmarshal(detail, &det) //convert json to struct
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		display(res, w)
		log.Error(err)
		return
	}
	//validation
	if det.Username == "" || det.Username == " " || det.Otp == "" || det.Otp == " " ||
		det.Newpassword == "" || det.Newpassword == " " {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter all the details"
		display(res, w)
		log.Warn(res["message"])
		return
	}
	if flag := validation.Passwordvalidation(res, det.Newpassword, w); flag {
		return
	}
	//Username exist or not
	err = db.Model(&det1).Where("username=?", det.Username).Select()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please enter valid Username "
		display(res, w)
		log.Warn(err)
		return
	}
	//checks otp matches with database
	err = bcrypt.CompareHashAndPassword([]byte(det1.Otp), []byte(det.Otp))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized) //status code for unathorization
		res["message"] = "OTP Entered is wrong!!!"
		display(res, w)
		log.Warn(err)
		return
	}

	//Encryption of password
	if bytepass = validation.Encrption(det.Newpassword, w, res); bytepass == nil {
		return
	}
	//updating password into database
	_, err = db.Model(&det1).Set("password=?", string(bytepass)).Where("username=?", det.Username).Update()
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		res["message"] = "Password reset failed"
		display(res, w)
		log.Error(err)
	} else {
		w.WriteHeader(http.StatusCreated)
		res["message"] = "Password reset success"
		display(res, w)
		log.Info(res["message"])
	}

}
