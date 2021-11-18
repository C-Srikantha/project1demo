package endpoints

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"project1.com/project/dbconnection"
)

func TestRegister(t *testing.T) {
	db, _ := dbconnection.DatabaseConnection()
	input := []byte(`{"firstname":"Srikantha","lastname":"c","username":"sriki","password":"aa@123AA","email":"srik@gmail.com"}`)
	req, err := http.NewRequest("POST", "/registration", bytes.NewBuffer(input))
	if err != nil {
		fmt.Println(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { PostRegistration(rw, r, db) })
	handler.ServeHTTP(rr, req)
	if statuscode := rr.Code; statuscode != http.StatusCreated {
		t.Errorf("got=%v,want=%v", rr.Code, http.StatusCreated)
	}

}
func TestLogin(t *testing.T) {
	db, _ := dbconnection.DatabaseConnection()
	input := []byte(`{"username":"sriki","password":"aa@123AA"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(input))
	if err != nil {
		fmt.Println(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { Login(rw, r, db) })
	handler.ServeHTTP(rr, req)
	if statuscode := rr.Code; statuscode != http.StatusFound {
		t.Errorf("got=%v,want=%v", rr.Code, http.StatusFound)
	}

}
