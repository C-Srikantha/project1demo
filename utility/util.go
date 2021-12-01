package utility

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//Encrption function encrpts the password and returns byte value and a error
func Encrption(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10) //Encryption of password feild to bytes
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
func Display(res map[string]string, w http.ResponseWriter) {
	jsonstr, _ := json.Marshal(res)
	w.Write(jsonstr)
}
