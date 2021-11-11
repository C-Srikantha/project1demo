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

func Login(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	var flag bool = true
	detail, err := io.ReadAll(r.Body)
	error(err)
	var det Logininfo
	var det1 Registration
	err = json.Unmarshal(detail, &det)
	error(err)
	err = db.Model(&det1).Where("username=?", det.Username).Select()
	if err != nil {
		fmt.Fprint(w, "No user found")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(det1.Password), []byte(det.Password))
	if err != nil {
		flag = false
	}
	if flag == true {
		fmt.Fprintf(w, "%s Welcome", det.Username)
	} else {
		fmt.Fprint(w, "Entered password is wrong")
	}

}
