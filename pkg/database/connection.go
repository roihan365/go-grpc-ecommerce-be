package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDB(ctx context.Context, connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		fmt.Println("Error when connecting to database:", err)
		panic(err)
	}

	return db
}
