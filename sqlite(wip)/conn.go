package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteDB struct {
	host, user, dbname, pass, charset string
	port                              int
	dbConn                            *sql.DB
	fieldlist                         []string
	tx                                *sql.Tx
}
