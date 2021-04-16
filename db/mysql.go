package db

import (
	"database/sql"
	"strings"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct{ Dsn string }

// initDB - DB파일 생성
func (d *Mysql) initDB() (*sql.DB, error) {
	var err error

	Dbo, err = sql.Open("mysql", d.Dsn)
	if err != nil {
		return nil, err
	}

	return Dbo, nil
}

// CreateTable - 테이블 생성
func (d *Mysql) CreateTable(recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		IDX INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
		NAME VARCHAR(128) NULL DEFAULT NULL COLLATE 'utf8_general_ci',
		PRICE DOUBLE NULL DEFAULT NULL,
		AUTHOR VARCHAR(128) NULL DEFAULT NULL COLLATE 'utf8_general_ci',
		ISBN VARCHAR(13) NULL DEFAULT NULL COLLATE 'utf8_general_ci',
		PRIMARY KEY (IDX),
		UNIQUE INDEX ISBN (ISBN),
		INDEX IDX (IDX)
	)
	COLLATE='utf8_general_ci'
	ENGINE=InnoDB;`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", TableName)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}