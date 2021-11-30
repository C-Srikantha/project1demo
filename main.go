package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"project1.com/project/createtable"
	"project1.com/project/dbconnection"
	"project1.com/project/endpoints"
)

//prints the errors to the user end
func error(err interface{}) {
	if err != nil {
		fmt.Println(err) //
		os.Exit(1)
	}
}

func handleRequest(db *pg.DB) {
	mux := mux.NewRouter()
	mux.HandleFunc("/registration", func(rw http.ResponseWriter, r *http.Request) { endpoints.PostRegistration(rw, r, db) }).Methods("POST")
	mux.HandleFunc("/login", func(rw http.ResponseWriter, r *http.Request) { endpoints.Login(rw, r, db) }).Methods("POST")
	mux.HandleFunc("/generateotp", func(rw http.ResponseWriter, r *http.Request) { endpoints.Resetpassotp(rw, r, db) }).Methods("POST")
	mux.HandleFunc("/reset", func(rw http.ResponseWriter, r *http.Request) { endpoints.Reset(rw, r, db) }).Methods("POST")
	err := http.ListenAndServe(":8081", mux)
	error(err)

}
func main() {
	db, err := dbconnection.DatabaseConnection() //calling database connection and returns db connection and error
	error(err)
	err = createtable.CreateTable(db) //calling createtable and returns error
	error(err)
	handleRequest(db)
}
