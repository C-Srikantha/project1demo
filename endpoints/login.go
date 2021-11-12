package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
)

type Logininfo struct {
	Username string
	Password string
}

//checks user is present and allows user to login is password matches in db
func Login(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	detail, err := io.ReadAll(r.Body)
	if err != nil {
		badrequesterror(err, w)
		return
	}
	var det Logininfo
	var det1 Registration
	err = json.Unmarshal(detail, &det)
	if err != nil {
		error(err)
		return
	}
	res := map[string]string{"message": ""}
	err = db.Model(&det1).Where("username=?", det.Username).Select()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res["message"] = "No User Found"
		str, err := json.Marshal(res)
		if err != nil {
			error(err)
			return
		}
		w.Write(str)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(det1.Password), []byte(det.Password)) //decrypt password
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized) //status code for unathorization
		res["message"] = "Entered password is wrong!!!"
		str, err := json.Marshal(res)
		if err != nil {
			error(err)
			return
		}
		w.Write(str)
		return
	} else {
		str := fmt.Sprintf("%s Welcome", det.Username)
		res["message"] = str
		jsonstr, err := json.Marshal(res)
		if err != nil {
			error(err)
			return
		}
		w.Write(jsonstr)
		return
	}

}
