package db // import "db"

import "database/sql"

const (
	_ = iota
	// MYSQL - mysql
	MYSQL = "mysql"
	// SQLITE - sqlite
	SQLITE = "sqlite3"
)

var (
	Filename    = "./data.db"
	TableName   = "BOOKSHELF"
	Dsn         *sql.DB
	UpdateScope []string     // UPDATE ... WHERE idx=?
	IgnoreScope []string     // Ignore if nil or null
	listCount   uint     = 3 // Default list count
	OrderScope  string       // Default order column name
)
