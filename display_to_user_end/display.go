package display_to_user_end

import (
	"encoding/json"
	"net/http"
)

func Display(res map[string]string, w http.ResponseWriter) {
	jsonstr, _ := json.Marshal(res)
	w.Write(jsonstr)
}
