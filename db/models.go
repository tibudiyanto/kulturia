// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql"
)

type Asset struct {
	ID       int64
	EntryID  sql.NullInt64
	Location sql.NullString
}

type Entry struct {
	ID     int64
	Name   string
	Origin string
	Desc   string
}
