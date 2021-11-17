package logsetup

import (
	"encoding/json"
	"net/http"
	"os"
)

//displays errors to user end
func error(res map[string]string, w http.ResponseWriter) {
	jsonstr, _ := json.Marshal(res)
	w.Write(jsonstr)
}
func Logfile(w http.ResponseWriter, res map[string]string) (*os.File, bool) {
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) //opening a log file
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res["message"] = "Something wrong in backend..Cant Open Log file"
		error(res, w)
		return nil, true
	}
	return file, false
}
