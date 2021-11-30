package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"project1.com/project/createtable"
	"project1.com/project/dbconnection"
	"project1.com/project/endpoints"
	"project1.com/project/logsetup"
)

var file *os.File

//prints the errors to the user end
func error(err interface{}) {
	log.SetOutput(file)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func handleRequest(db *pg.DB, file *os.File) {
	mux := mux.NewRouter()
	mux.HandleFunc("/registration", func(rw http.ResponseWriter, r *http.Request) { endpoints.PostRegistration(rw, r, db, file) }).Methods("POST")
	mux.HandleFunc("/login", func(rw http.ResponseWriter, r *http.Request) { endpoints.Login(rw, r, db, file) }).Methods("POST")
	mux.HandleFunc("/generateotp", func(rw http.ResponseWriter, r *http.Request) { endpoints.ResetPassotp(rw, r, db, file) }).Methods("POST")
	mux.HandleFunc("/reset", func(rw http.ResponseWriter, r *http.Request) { endpoints.Reset(rw, r, db, file) }).Methods("POST")
	err := http.ListenAndServe(":8081", mux)
	error(err)

}
func main() {
	file, err := logsetup.LogFile()
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	db, err := dbconnection.DatabaseConnection() //calling database connection and returns db connection and error
	error(err)
	err = createtable.CreateTable(db) //calling createtable and returns error
	error(err)
	handleRequest(db, file)
}
