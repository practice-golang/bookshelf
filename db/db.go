package db // import "db"

import (
	"database/sql"
	"errors"
	"sync"
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
		CreateTable(bool) error
	}
)

var (
	Dbi          DBI    // DB Object Interface
	Dsn          string // Data Source Name
	once         sync.Once
	DatabaseName = "bookshelf"
	TableName    = "books"
	Dbo          *sql.DB
	DBType       int
	UpdateScope  []string     // UPDATE ... WHERE idx=?
	IgnoreScope  []string     // Ignore if nil or null
	listCount    uint     = 3 // Default list count
	OrderScope   string       // Default order column name
)

// InitDB - DB 준비
func InitDB(driver int) (DBI, error) {
	var err error
	var dbo DBI
	once.Do(func() {
		dbo, err = dbFactory(driver)
	})
	return dbo, err
}

func dbFactory(driver int) (DBI, error) {
	switch driver {
	case MYSQL:
		dbo := &Mysql{Dsn: Dsn}
		dbo.initDB()
		return dbo, nil
	case SQLSERVER:
		dbo := &Sqlserver{Dsn: Dsn}
		dbo.initDB()
		return dbo, nil
	case SQLITE:
		dbo := &Sqlite{Dsn: Dsn}
		dbo.initDB()
		return dbo, nil
	default:
		return nil, errors.New("nothing to support DB")
	}
}
