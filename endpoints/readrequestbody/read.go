package readrequestbody

import (
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	user "project1.com/project/display_to_user_end"
)

//reads the body from request and converts json to struct ...returns err
func Readbody(r *http.Request, w http.ResponseWriter, res map[string]string, det interface{}) error {
	detail, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		user.Display(res, w)
		log.Error(err)
		return err
	}
	err = json.Unmarshal(detail, &det) //converts json format to struct
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		user.Display(res, w)
		log.Error(err)
		return err
	}
	return nil
}
