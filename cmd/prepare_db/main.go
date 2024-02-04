package main

import (
	"context"
	"database/sql"
	"fmt"
	"kulturia/db"
	"kulturia/sqlc"

	_ "github.com/mattn/go-sqlite3"
)

// Move this somewhere else or as another executable
func prepareDatabase(dbx db.DBTX) {
	ctx := context.Background()
	if _, err := dbx.ExecContext(ctx, sqlc.DDL); err != nil {
		fmt.Println("err", err)
	}
}

func main() {
	dbx, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}
	prepareDatabase(dbx)
}
