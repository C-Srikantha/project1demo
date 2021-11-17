package validation

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//encrption of password
func Encrption(password string, w http.ResponseWriter, res map[string]string) []byte {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10) //Encryption of password feild to bytes
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend...Failed to encrypt password"
		error(res, w)
		log.Print(err)
		return nil
	}
	return bytes
}
