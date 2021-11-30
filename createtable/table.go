package createtable

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

//defines the tabel feilds using struct
type Registration struct {
	Id        int `pg:",pk"`
	Firstname string
	Lastname  string
	Username  string `sql:",unique"`
	Password  string
	Email     string `sql:",unique"`
	Otp       string
}

//Creates tabel in database and returns error
func CreateTable(db *pg.DB) error {
	table := []interface{}{
		(*Registration)(nil),
	}
	for _, table := range table {
		err := db.Model(table).CreateTable(&orm.CreateTableOptions{ //query for creating tabel
			IfNotExists: true, //checks if table exists or not
			Varchar:     100,  //converts all string type to varchar in db
		})
		if err != nil {
			return err
		}
	}
	return nil
}
