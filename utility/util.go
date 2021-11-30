package utility

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"project1.com/project/logsetup"
)

//encrption of password
func Encrption(password string, w http.ResponseWriter, res map[string]string) []byte {
	file, flag := logsetup.Logfile(w, res)
	defer file.Close()
	if flag {
		return nil
	}
	log.SetOutput(file)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10) //Encryption of password feild to bytes
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend...Failed to encrypt password"
		Display(res, w)
		log.Println(err.Error()) //log at endpoints
		return nil
	}
	return bytes
}
func Display(res map[string]string, w http.ResponseWriter) {
	jsonstr, _ := json.Marshal(res)
	w.Write(jsonstr)
}
