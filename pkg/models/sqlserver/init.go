package sqlserver

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

// Returns an initialized sql.DB connection and error. First forms a connection string
// from a .env, then opens a connection and performs a ping to ensure that a valid conneciton
// has been established.
func InitDB() (*sql.DB, error) {
	UN, PW, SRV, DB := os.Getenv("UN"), os.Getenv("PW"), os.Getenv("SRV"), os.Getenv("DB")
	PORT := 1433
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		SRV, UN, PW, PORT, DB)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, errors.New("Database error: " + err.Error())
	}
	if err = db.Ping(); err != nil {
		return nil, errors.New("Error connecting to db: " + err.Error())
	}
	log.Println("Connection opened")
	return db, nil
}
