package readrequestbody

import (
	"encoding/json"
	"io"
	"net/http"
)

//Readbody function reads request body and returns byte value and error
func ReadBody(r *http.Request) ([]byte, error) {
	detail, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

//Covert function converts json to struct
func Convert(detail []byte, det interface{}) error {
	err := json.Unmarshal(detail, &det) //converts json format to struct
	if err != nil {
		return err
	}
	return nil
}
