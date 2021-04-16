package db

import (
	"database/sql"
	"strings"
	// _ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
)

type Sqlite struct{ Dsn string }

// initDB - DB파일 생성
func (d *Sqlite) initDB() (*sql.DB, error) {
	var err error

	Dbo, err = sql.Open("sqlite", d.Dsn)
	if err != nil {
		return nil, err
	}

	return Dbo, nil
}

// CreateTable - 테이블 생성
func (d *Sqlite) CreateTable(recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS "#TABLE_NAME";`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS "#TABLE_NAME" (
		"IDX"			INTEGER,
		"NAME"			TEXT,
		"PRICE"			REAL,
		"AUTHOR"		TEXT,
		"ISBN"			TEXT UNIQUE,
		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", TableName)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
