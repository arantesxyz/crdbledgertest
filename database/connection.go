package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetConnection(uri string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
