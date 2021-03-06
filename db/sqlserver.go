package db

import (
	"database/sql"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
)

type Sqlserver struct{ Dsn string }

// initDB - Prepare DB
func (d *Sqlserver) initDB() (*sql.DB, error) {
	var err error

	Dbo, err = sql.Open("sqlserver", d.Dsn)
	if err != nil {
		return nil, err
	}

	return Dbo, nil
}

func (d *Sqlserver) CreateDB() error {
	return nil
}

// CreateTable - Create table
func (d *Sqlserver) CreateTable(recreate bool) error {
	sql := `
	USE master
	-- GO

	IF NOT EXISTS(
		SELECT name
		FROM sys.databases
		WHERE name=N'#DATABASE'
	)
	CREATE DATABASE [#DATABASE]
	-- GO
	`
	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	if recreate {
		sql = `USE #DATABASE`
		sql += `
		IF OBJECT_ID('#TABLE_NAME','U') IS NOT NULL
		DROP TABLE #TABLE_NAME
		-- GO
		`

		sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
		sql = strings.ReplaceAll(sql, "#TABLE_NAME", TableName)

		_, err := Dbo.Exec(sql)
		if err != nil {
			return err
		}
	}

	sql = `USE #DATABASE`
	sql += `
	IF OBJECT_ID(N'#TABLE_NAME', N'U') IS NULL
	CREATE TABLE #TABLE_NAME (
		IDX INT NOT NULL IDENTITY PRIMARY KEY,
		NAME VARCHAR(128) NOT NULL,
		PRICE DECIMAL(10,2) NOT NULL,
		AUTHOR VARCHAR(128) NOT NULL,
		ISBN VARCHAR(128) NOT NULL UNIQUE,
	)
	--GO`

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", TableName)

	_, err = Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
