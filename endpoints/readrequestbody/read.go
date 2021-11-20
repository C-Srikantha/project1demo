package readrequestbody

import (
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"project1.com/project/utility"
)

//reads the body from request and converts json to struct ...returns err
func Readbody(r *http.Request, w http.ResponseWriter, res map[string]string, det interface{}) error {
	detail, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		utility.Display(res, w)
		log.Error(err)
		return err
	}
	err = json.Unmarshal(detail, &det) //converts json format to struct
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant convert json to struct"
		utility.Display(res, w)
		log.Error(err)
		return err
	}
	return nil
}
