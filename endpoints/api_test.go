package endpoints

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"project1.com/project/dbconnection"
)

type Inputs struct {
	input    []byte
	wantcode int
}

var reg = []Inputs{
	{[]byte(`{"firstname":"Srikantha","lastname":"c","username":"sriki","password":"aa@123AA","email":"srikan@getnada.com"}`), http.StatusCreated},
	{[]byte(`{"firstname":"Srikantha","lastname":"D","username":"sriki","password":"aa@123AA","email":"srika@getnada.com"}`), http.StatusAlreadyReported},
	{[]byte(`{"firstname":"","lastname":"","username":"sriki123","password":"aa@123AA","email":""}`), http.StatusBadRequest},
	{[]byte(`{"firstname":"Hello","lastname":"h","username":"hello","password":"aa@123AA","email":"srikan@getnada.com"}`), http.StatusAlreadyReported},
	{[]byte(`{"firstname":"Hello","lastname":"l","username":"hello123","password":"a123AA","email":"siwis@givmail.com"}`), http.StatusBadRequest},
	{[]byte(`{"firstname":"Hello","lastname":"l","username":"hello321","password":"aa@123AA","email":"siwisgivmailcom"}`), http.StatusBadRequest},
}
var login = []Inputs{
	{[]byte(`{"username":"sriki","password":"aa@123AA"}`), http.StatusFound},
	{[]byte(`{"username":"","password":"aa@123AA"}`), http.StatusBadRequest},
	{[]byte(`{"username":"sriki","password":""}`), http.StatusBadRequest},
	{[]byte(`{"username":"hello1234","password":"aa@123AA"}`), http.StatusNotFound},
	{[]byte(`{"username":"sriki","password":"a@123AA"}`), http.StatusUnauthorized},
}

func TestRegister(t *testing.T) {
	db, _ := dbconnection.DatabaseConnection()
	for _, val := range reg {
		req, err := http.NewRequest("POST", "/registration", bytes.NewBuffer(val.input))
		if err != nil {
			fmt.Println(err)
			return
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { PostRegistration(rw, r, db) })
		handler.ServeHTTP(rr, req)
		if statuscode := rr.Code; statuscode != val.wantcode {
			t.Errorf("got=%v,want=%v", rr.Code, val.wantcode)
		}
	}

}
func TestLogin(t *testing.T) {
	db, _ := dbconnection.DatabaseConnection()
	for _, val := range login {
		//input := []byte(`{"username":"sriki","password":"aa@123AA"}`)
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(val.input))
		if err != nil {
			fmt.Println(err)
			return
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { Login(rw, r, db) })
		handler.ServeHTTP(rr, req)
		if statuscode := rr.Code; statuscode != val.wantcode {
			t.Errorf("got=%v,want=%v", rr.Code, val.wantcode)
		}
	}

}
