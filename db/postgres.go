package db

import (
	"database/sql"
	"strings"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/lib/pq"
)

type Postgres struct{ Dsn string }

// initDB - Prepare DB
func (d *Postgres) initDB() (*sql.DB, error) {
	var err error

	Dbo, err = sql.Open("postgres", d.Dsn)
	if err != nil {
		return nil, err
	}

	return Dbo, nil
}

// CreateTable - Create table
func (d *Postgres) CreateTable(recreate bool) error {
	sql := `CREATE SCHEMA IF NOT EXISTS #SCHEMA_NAME;`

	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		"IDX" SERIAL PRIMARY KEY,
		"NAME" VARCHAR(128) NULL DEFAULT NULL,
		"PRICE" NUMERIC(10,2) NULL DEFAULT NULL,
		"AUTHOR" VARCHAR(128) NULL DEFAULT NULL,
		"ISBN" VARCHAR(13) UNIQUE NULL DEFAULT NULL
	);`

	sql = strings.ReplaceAll(sql, "#SCHEMA_NAME", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", TableName)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
