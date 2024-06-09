package db

import (
	"fmt"
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type sqlStr struct {
	driver string
	conStr string
}

func (s *sqlStr) makeConnection() (*sql.DB, error) {
	connection, err := sql.Open(s.driver, s.conStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	return connection, nil
}

/*
* e.g conn, err := db.New("postgres", "postgres://postgres:password@localhost:port_no/db_name")
*/ 
func New(driverName string, conStr string) (*sql.DB, error) {
	c := &sqlStr{conStr: conStr, driver: driverName}
	return c.makeConnection()
}
