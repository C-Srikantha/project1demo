package main

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"project1.com/project/createtable"
	"project1.com/project/dbconnection"
)

//prints the errors to the user end
func error(err interface{}) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func handlerequest(db *pg.DB) {

}
func main() {
	db, err := dbconnection.DatabaseConnection() //calling database connection and returns db connection and error
	error(err)
	err = createtable.Createtable(db) //calling createtable and returns error
	error(err)
	handlerequest(db)
}
