package db

import (
	"errors"
	"log"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	"github.com/thoas/go-funk"
	"gopkg.in/guregu/null.v4"
)

// PrepareWhere - Where 요소 준비
func PrepareWhere(data interface{}) goqu.Ex {
	result := goqu.Ex{}

	values := reflect.ValueOf(data)
	dataReflect := values.Type()

	for i := 0; i < values.NumField(); i++ {
		f := values.Field(i).Interface()
		b := dataReflect.Field(i)
		jsonName := b.Tag.Get("json")
		switch f := f.(type) {
		case null.String:
			if f.Valid {
				result[jsonName], _ = f.Value()
			}
		case null.Int:
			if f.Valid {
				result[jsonName], _ = f.Value()
			}
		case null.Float:
			if f.Valid {
				result[jsonName], _ = f.Value()
			}
		}
	}

	log.Println(result)

	return result
}

// CheckValidAndPrepareWhere - 빠진 요소 체크, Where 요소 준비
func CheckValidAndPrepareWhere(book interface{}) (goqu.Ex, error) {
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
			if funk.Contains(UpdateScope, jsonName) {
				result[jsonName], _ = f.Value()
			}
		case null.Int:
			if !f.Valid {
				return nil, errors.New("`" + jsonName + "` must have a value")
			}
			if funk.Contains(UpdateScope, jsonName) {
				result[jsonName], _ = f.Value()
			}
		case null.Float:
			if !f.Valid {
				return nil, errors.New("`" + jsonName + "` must have a value")
			}
			if funk.Contains(UpdateScope, jsonName) {
				result[jsonName], _ = f.Value()
			}
		}
	}

	return result, nil
}
