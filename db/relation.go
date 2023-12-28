package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type sql struct {
	conStr string
}

func (s *sql) makeConnection() *pgx.Conn {
	connection, err := pgx.Connect(context.Background(), s.conStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// defer connection.Close(context.Background())

	return connection
}

// urlExample := "postgres://username:password@localhost:5432/database_name"
func New(conStr string) *pgx.Conn {
	c := &sql{ conStr: conStr }
	return c.makeConnection()
}