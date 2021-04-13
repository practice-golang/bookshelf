package db

import (
	"database/sql"
	"errors"
	"log"
	"models"
	"reflect"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/thoas/go-funk"
	"gopkg.in/guregu/null.v4"
)

// Sqlite - sqlite
type Sqlite struct {
	TableName string
}

// CreateTable - 테이블 생성
// func (s *Sqlite) CreateTable(recreate bool) error {
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

// CheckValidAndPrepareWhere - 빠진 요소 체크, Where 요소 준비
func CheckValidAndPrepareWhere(book models.Book) (goqu.Ex, error) {
	result := goqu.Ex{}

	values := reflect.ValueOf(book)
	bookReflect := values.Type()

	for i := 0; i < values.NumField(); i++ {
		f := values.Field(i).Interface()
		b := bookReflect.Field(i)
		// log.Println(b.Name)
		jsonName := b.Tag.Get("json")
		switch f := f.(type) {
		case null.String:
			if !f.Valid {
				return nil, errors.New("`" + jsonName + "` must have a value")
			}
			if funk.Contains(UpdateTarget, jsonName) {
				result[jsonName], _ = f.Value()
			}
		case null.Int:
			if !f.Valid {
				return nil, errors.New("`" + jsonName + "` must have a value")
			}
			if funk.Contains(UpdateTarget, jsonName) {
				result[jsonName], _ = f.Value()
			}
		case null.Float:
			if !f.Valid {
				return nil, errors.New("`" + jsonName + "` must have a value")
			}
			if funk.Contains(UpdateTarget, jsonName) {
				result[jsonName], _ = f.Value()
			}
		}
	}

	return result, nil
}

// InsertData - Crud
func InsertData(books []models.Book) (sql.Result, error) {
	dbms := goqu.New("sqlite3", Dsn)
	ds := dbms.Insert(TableName).Rows(books)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dsn.Exec(sql)
	if err != nil {
		// lastID, _ := result.LastInsertId()
		// affRows, _ := result.RowsAffected()
		// log.Println(lastID, affRows)
		return nil, err
	}

	return result, nil
}

// SelectData - cRud
func SelectData(search models.Book) ([]models.Book, error) {
	var result []models.Book

	dbms := goqu.New("sqlite3", Dsn)
	ds := dbms.From(TableName).Select(&models.Book{})
	err := ds.ScanStructs(&result)
	if err != nil {
		log.Println("ds: ", err.Error())
		return nil, err
	}

	return result, nil
}

// UpdateData - crUd
func UpdateData(book models.Book) (sql.Result, error) {
	whereEXP, err := CheckValidAndPrepareWhere(book)
	if err != nil {
		return nil, err
	}

	dbms := goqu.New("sqlite3", Dsn)
	ds := dbms.Update(TableName).Set(book).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dsn.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteData(target, value string) (sql.Result, error) {
	whereEXP := goqu.Ex{target: value}

	dbms := goqu.New("sqlite3", Dsn)
	ds := dbms.Delete(TableName).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dsn.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}
