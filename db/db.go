package db

import (
	"database/sql"
	"errors"
)

const (
	_ = iota
	// SQLITE - sqlite
	SQLITE
	// SQLSERVER - mssql
	SQLSERVER
	// MYSQL - mysql
	MYSQL
	// POSTGRES - pgsql
	POSTGRES
)

type (
	DBI interface {
		initDB() (*sql.DB, error)
		CreateDB() error
		CreateTable(bool) error
	}
)

var (
	Dbi DBI    // DB Object Interface
	Dsn string // Data Source Name
	// once         sync.Once
	DatabaseName = "bookshelf"
	TableName    = "books"
	Dbo          *sql.DB
	DBType       int
	UpdateScope  []string     // UPDATE ... WHERE IDX=?
	IgnoreScope  []string     // Ignore if nil or null
	listCount    uint     = 3 // Default list count
	OrderScope   string       // Default order column name
)

// InitDB - Prepare DB
func InitDB(driver int) (DBI, error) {
	var err error
	var dbi DBI

	dbi, err = dbFactory(driver)

	return dbi, err
}

func dbFactory(driver int) (DBI, error) {
	switch driver {
	case SQLITE:
		dbi := &Sqlite{Dsn: Dsn}
		dbi.initDB()
		return dbi, nil
	case MYSQL:
		dbi := &Mysql{Dsn: Dsn}
		dbi.initDB()
		return dbi, nil
	case SQLSERVER:
		dbi := &Sqlserver{Dsn: Dsn}
		dbi.initDB()
		return dbi, nil
	case POSTGRES:
		dbi := &Postgres{Dsn: Dsn}
		dbi.initDB()
		return dbi, nil
	default:
		return nil, errors.New("nothing to support DB")
	}
}
