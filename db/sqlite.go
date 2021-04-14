package db

import (
	"database/sql"
	"strings"
)

// InitDB - DB파일 생성
func InitDB() (*sql.DB, error) {
	var err error
	fn := Filename

	Dsn, err = sql.Open("sqlite", fn)
	if err != nil {
		return nil, err
	}

	return Dsn, nil
}

// CreateTable - 테이블 생성
func CreateTable(recreate bool) error {
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

	_, err := Dsn.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
