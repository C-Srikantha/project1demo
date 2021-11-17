package endpoints

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-pg/pg"
	"project1.com/project/logsetup"
	"project1.com/project/otp"
	"project1.com/project/validation"
)

type Resetpassword struct {
	Username string
}

func Resetpassotp(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	file, flag := logsetup.Logfile(w, res)
	defer file.Close()
	if flag {
		return
	}
	//reads username from body
	detail, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		error(res, w)
		log.Println(err.Error())
		return
	}
	var det Resetpassword
	var det1 Registration
	err = json.Unmarshal(detail, &det) //convert json to struct
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		error(res, w)
		log.Println(err.Error())
		return
	}
	//validation
	if det.Username == "" || det.Username == " " {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter all the details"
		error(res, w)
		return
	}
	//checks username exists or not
	err = db.Model(&det1).Where("username=?", det.Username).Select()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res["message"] = "User Not Found"
		error(res, w)
		log.Print(err.Error())
		return
	}
	//calling generate otp func
	otp, flag := otp.Generateotp(det1.Email)
	if flag {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Failed to send mail"
		error(res, w)
	}
	//encryption of otp
	if bytes = validation.Encrption(otp, w, res); bytes == nil {
		return
	}
	//updating otp feild in database
	_, err = db.Model(&det1).Set("otp=?", string(bytes)).Where("username=?", det.Username).Update()
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		res["message"] = "Something wrong in backend"
		error(res, w)
		log.Println(err.Error())
	} else {
		w.WriteHeader(http.StatusCreated)
		res["message"] = "Otp has sent via mail"
		error(res, w)
	}
}
