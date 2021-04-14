package db // import "db"

import "database/sql"

var (
	Filename    = "./data.db"
	TableName   = "BOOKSHELF"
	Dsn         *sql.DB
	UpdateScope []string // UPDATE ... WHERE idx=?
	IgnoreScope []string // Ignore if nil or null
	listCount   uint     = 3
)

const (
	_ = iota
	// MYSQL - mysql
	MYSQL = "mysql"
	// SQLITE - sqlite
	SQLITE = "sqlite3"
)
