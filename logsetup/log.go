package logsetup

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	users "project1.com/project/display_to_user_end"
)

//displays errors to user end
/*func error(res map[string]string, w http.ResponseWriter) {
	jsonstr, _ := json.Marshal(res)
	w.Write(jsonstr)
}*/
func Logfile(w http.ResponseWriter, res map[string]string) (*os.File, bool) {
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) //opening a log file
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant Open Log file"
		users.Display(res, w)
		return nil, true
	}
	formater := new(log.TextFormatter)
	formater.TimestampFormat = "02-01-2006 15:04:05"
	formater.FullTimestamp = true
	log.SetFormatter(formater)

	return file, false
}
