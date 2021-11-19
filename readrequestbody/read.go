package readrequestbody

import (
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	user "project1.com/project/display_to_user_end"
)

func Readbody(r *http.Request, w http.ResponseWriter, res map[string]string) ([]byte, error) {
	detail, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res["message"] = "Failed to read request body!!!"
		user.Display(res, w)
		log.Error(err)
		return nil, err
	}
	return detail, nil
}
