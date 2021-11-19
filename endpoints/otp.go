package endpoints

import (
	"net/http"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	read "project1.com/project/endpoints/readrequestbody"
	"project1.com/project/logsetup"
	"project1.com/project/otp"
	"project1.com/project/validation"
)

type Resetpassword struct {
	Username string `validate:"nonzero"`
}

func Resetpassotp(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	file, flag := logsetup.Logfile(w, res)
	if flag {
		return
	}
	defer file.Close()
	log.SetOutput(file)
	var det *Resetpassword
	var det1 Registration
	//reads username from body
	if err := read.Readbody(r, w, res, &det); err != nil {
		return
	}
	/*detail, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		display(res, w)
		log.Error(err)
		return
	}
	var det Resetpassword
	var det1 Registration
	err = json.Unmarshal(detail, &det) //convert json to struct
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		display(res, w)
		log.Error(err)
		return
	}*/
	//validation
	/*if err := validator.Validate(det); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter all the details"
		display(res, w)
		log.Warn(res["message"])
		return
	}*/
	//validation
	if err := validation.FeildValidation(det, w, res); err != nil {
		return
	}
	//checks username exists or not
	err := db.Model(&det1).Where("username=?", det.Username).Select()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res["message"] = "User Not Found"
		display(res, w)
		log.Warn(err)
		return
	}
	//calling generate otp func
	otpstr, flag := otp.Generateotp()
	if flag {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "generating otp failed"
		display(res, w)
		log.Error(otpstr)
		return
	}
	//encryption of otp
	if bytepass = validation.Encrption(otpstr, w, res); bytepass == nil {
		return
	}
	//updating otp feild in database
	_, err = db.Model(&det1).Set("otp=?", string(bytepass)).Where("username=?", det.Username).Update()
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		res["message"] = "Something wrong in backend..Failed to update pass"
		display(res, w)
		log.Error(err)
	} else {
		w.WriteHeader(http.StatusCreated)
		if err := otp.Emailgenerate(det1.Email, otpstr); err != nil {
			w.WriteHeader(http.StatusNotModified)
			res["message"] = "Failed to send mail"
			display(res, w)
			log.Error(err)
			return
		}
		res["message"] = "Otp has sent via mail"
		display(res, w)
		log.Info(res["message"])
	}

}
