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

func (s *sqlStr) makeConnection() *sql.DB {
	connection, err := sql.Open(s.driver, s.conStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return connection
}

func New(driverName string, conStr string) *sql.DB {
	c := &sqlStr{ conStr: conStr, driver: driverName }
	return c.makeConnection()
}