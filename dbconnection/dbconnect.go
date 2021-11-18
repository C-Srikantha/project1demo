package dbconnection

import (
	"context"

	"github.com/go-pg/pg"
)

//connection to the database
func DatabaseConnection() (*pg.DB, error) {
	dbdetail := &pg.Options{
		User:     "postgres",
		Password: "codecraft",
		Addr:     ":8080",
		Database: "project1demo",
	}
	connection := pg.Connect(dbdetail)
	control := context.Background()
	_, err := connection.ExecContext(control, "SELECT 1")
	return connection, err

}
