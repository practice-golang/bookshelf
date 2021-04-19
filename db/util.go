package db

import (
	"errors"
	"log"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	"github.com/thoas/go-funk"
	"gopkg.in/guregu/null.v4"
)

// PrepareWhere - Prepare where condition for search
func PrepareWhere(data interface{}) goqu.Ex {
	result := goqu.Ex{}

	values := reflect.ValueOf(data)
	dataReflect := values.Type()

	for i := 0; i < values.NumField(); i++ {
		f := values.Field(i).Interface()
		b := dataReflect.Field(i)
		// dbName := b.Tag.Get("json")
		dbName := b.Tag.Get("db")
		switch f := f.(type) {
		case null.String:
			if f.Valid {
				result[dbName], _ = f.Value()
			}
		case null.Int:
			if f.Valid {
				result[dbName], _ = f.Value()
			}
		case null.Float:
			if f.Valid {
				result[dbName], _ = f.Value()
			}
		}
	}

	log.Println(result)

	return result
}

// CheckValidAndPrepareWhere - Prepare where and check missing values for update
func CheckValidAndPrepareWhere(book interface{}) (goqu.Ex, error) {
	result := goqu.Ex{}

	values := reflect.ValueOf(book)
	bookReflect := values.Type()

	for i := 0; i < values.NumField(); i++ {
		f := values.Field(i).Interface()
		b := bookReflect.Field(i)
		// log.Println(b.Name)
		// dbName := b.Tag.Get("json")
		dbName := b.Tag.Get("db")
		switch f := f.(type) {
		case null.String:
			if !f.Valid {
				return nil, errors.New("`" + dbName + "` must have a value")
			}
			if funk.Contains(UpdateScope, dbName) {
				result[dbName], _ = f.Value()
			}
		case null.Int:
			if !f.Valid {
				return nil, errors.New("`" + dbName + "` must have a value")
			}
			if funk.Contains(UpdateScope, dbName) {
				result[dbName], _ = f.Value()
			}
		case null.Float:
			if !f.Valid {
				return nil, errors.New("`" + dbName + "` must have a value")
			}
			if funk.Contains(UpdateScope, dbName) {
				result[dbName], _ = f.Value()
			}
		}
	}

	return result, nil
}
