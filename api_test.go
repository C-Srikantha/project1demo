package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"project1.com/project/dbconnection"
	"project1.com/project/endpoints"
	"project1.com/project/logsetup"
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
var inotp = []Inputs{
	{[]byte(`{"username":""}`), http.StatusBadRequest},
	{[]byte(`{"username":"srikantha"}`), http.StatusNotFound},
}
var reset = []Inputs{
	{[]byte(`{"username":"","otp":"","newpassword":""}`), http.StatusBadRequest},
	{[]byte(`{"username":"sriki","otp":"","newpassword":"sa@123AA"}`), http.StatusBadRequest},
	{[]byte(`{"username":"sriki","otp":"12345","newpassword":"sa@123AA"}`), http.StatusUnauthorized},
	{[]byte(`{"username":"srikiantha","otp":"12345","newpassword":"sa@123AA"}`), http.StatusBadRequest},
}

//test for regester api
func TestRegister(t *testing.T) {
	db, _ := dbconnection.DatabaseConnection()
	file, _ := logsetup.LogFile()
	defer file.Close()
	log.SetOutput(file)
	for _, val := range reg {
		req, err := http.NewRequest("POST", "/registration", bytes.NewBuffer(val.input))
		if err != nil {
			log.Error(err)
			return
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { endpoints.PostRegistration(rw, r, db, file) })
		handler.ServeHTTP(rr, req)
		if statuscode := rr.Code; statuscode != val.wantcode {
			t.Errorf("got=%v,want=%v", rr.Code, val.wantcode)
		}
	}

}

//test for login api
func TestLogin(t *testing.T) {
	db, _ := dbconnection.DatabaseConnection()
	file, _ := logsetup.LogFile()
	defer file.Close()
	log.SetOutput(file)
	for _, val := range login {
		//input := []byte(`{"username":"sriki","password":"aa@123AA"}`)
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(val.input))
		if err != nil {
			log.Error(err)
			return
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { endpoints.Login(rw, r, db, file) })
		handler.ServeHTTP(rr, req)
		if statuscode := rr.Code; statuscode != val.wantcode {
			t.Errorf("got=%v,want=%v", rr.Code, val.wantcode)
		}
	}

}

//test for otp api
func TestOtp(t *testing.T) {
	db, _ := dbconnection.DatabaseConnection()
	file, _ := logsetup.LogFile()
	defer file.Close()
	log.SetOutput(file)
	for _, val := range inotp {
		req, err := http.NewRequest("POST", "/generateotp", bytes.NewBuffer(val.input))
		if err != nil {
			log.Error(err)
			return
		}
		rr := httptest.NewRecorder()
		handle := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { endpoints.ResetPassotp(rw, r, db, file) })
		handle.ServeHTTP(rr, req)
		if statuscode := rr.Code; statuscode != val.wantcode {
			t.Errorf("got=%v,want=%v", rr.Code, val.wantcode)
		}
	}
}

//test for resetpass api
func TestReset(t *testing.T) {
	db, _ := dbconnection.DatabaseConnection()
	file, _ := logsetup.LogFile()
	defer file.Close()
	log.SetOutput(file)
	for _, val := range reset {
		req, err := http.NewRequest("POST", "/reset", bytes.NewBuffer(val.input))
		if err != nil {
			log.Error(err)
			return
		}
		rr := httptest.NewRecorder()
		handle := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { endpoints.Reset(rw, r, db, file) })
		handle.ServeHTTP(rr, req)
		if statuscode := rr.Code; statuscode != val.wantcode {
			t.Errorf("got=%v,want=%v", rr.Code, val.wantcode)
		}
	}
}
