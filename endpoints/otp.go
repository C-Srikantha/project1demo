package endpoints

import (
	"net/http"
	"os"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"

	"project1.com/project/createtable"
	read "project1.com/project/endpoints/readrequestbody"
	"project1.com/project/otp"
	"project1.com/project/utility"
	"project1.com/project/validation"
)

type Resetpassword struct { //naming convention
	Username string `validate:"nonzero"`
}

//ResetPassotp generates otp and sends mail to the user
func ResetPassotp(w http.ResponseWriter, r *http.Request, db *pg.DB, file *os.File) {
	log.SetOutput(file)
	var det *Resetpassword
	var det1 createtable.Registration
	//reads username from body
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
		log.Error(res["message"])
		return
	}
	//checks username exists or not
	err = db.Model(&det1).Where("username=?", det.Username).Select()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res["message"] = "User Not Found"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	//calling generate otp func
	otpstr, flag := otp.GenerateOtp()
	if flag {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "generating otp failed"
		utility.Display(res, w)
		log.Error(otpstr)
		return
	}
	//encryption of otp
	if bytepass, err = utility.Encrption(otpstr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend...Failed to encrypt password"
		utility.Display(res, w)
		log.Error(err)
		return
	}
	//updating otp feild in database
	_, err = db.Model(&det1).Set("otp=?", string(bytepass)).Where("username=?", det.Username).Update()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Failed to update pass"
		utility.Display(res, w)
		log.Error(err)
	} else {
		w.WriteHeader(http.StatusCreated)
		if err := otp.EmailGenerate(det1.Email, otpstr); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res["message"] = "Failed to send mail"
			utility.Display(res, w)
			log.Error(err)
			return
		}
		res["message"] = "Otp has sent via mail"
		utility.Display(res, w)
		log.Info(res["message"])
	}

}
