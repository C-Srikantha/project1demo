package endpoints

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-pg/pg"
	"project1.com/project/passemailvalidation"
)

type Resetpass struct {
	Username    string
	Oldpassword string
	Newpassword string
}

func Reset(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	file, flag := logfile(w)
	if flag {
		return
	}
	log.SetOutput(file) //setting output destination
	detail, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		error(res, w)
		log.Println(err)
		return
	}
	var det Resetpass
	var det1 Registration
	err = json.Unmarshal(detail, &det) //convert json to struct
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		error(res, w)
		log.Println(err)
		return
	}
	//validation
	if det.Username == "" || det.Username == " " || det.Oldpassword == "" || det.Oldpassword == " " ||
		det.Newpassword == "" || det.Newpassword == " " {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please Enter all the details"
		error(res, w)
		log.Println(err)
		return
	}
	if flag := passemailvalidation.Passwordvalidation(res, det.Newpassword, w); flag {
		return
	}
	//Username exist or not
	err = db.Model(&det1).Where("username=?", det.Username).Select()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Please enter valid Username and Pass"
		error(res, w)
		log.Println(err)
		return
	}
	//Checking Oldpassword is matches with database
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Oldpassword is wrong!!!"
		error(res, w)
		log.Println(err)
		return
	}
	//updating password into database
	_, err = db.Model(&det1).Set("password=?", det.Newpassword).Where("username=?", det.Username).Update()
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		res["message"] = "Password reset failed"
		error(res, w)
		log.Println(err)
	} else {
		w.WriteHeader(http.StatusCreated)
		res["message"] = "Password reset success"
		error(res, w)
		log.Println(err)

	}

}
