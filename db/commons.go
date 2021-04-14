package db

import (
	"database/sql"
	"log"
	"models"

	"github.com/doug-martin/goqu/v9"
)

// InsertData - Crud
func InsertData(data interface{}) (sql.Result, error) {
	dbms := goqu.New("sqlite3", Dsn)
	ds := dbms.Insert(TableName).Rows(data)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dsn.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SelectData - cRud
func SelectData(search interface{}) (interface{}, error) {
	var result interface{}
	var err error

	dbms := goqu.New("sqlite3", Dsn)
	ds := dbms.From(TableName).Select(search)

	bookResult := []models.Book{}

	switch search := search.(type) {
	case models.Book:
		whereEXP := PrepareWhere(search)
		if !whereEXP.IsEmpty() {
			ds = ds.Where(whereEXP)
		}
		// cnt := listCount
		cnt := uint(10)
		ds = ds.Limit(cnt)

	case models.BookSearch:
		keywords := search.Keywords
		whereEXPRs := []goqu.Expression{}
		for _, k := range keywords {
			exp := PrepareWhere(k)
			if !exp.IsEmpty() {
				exp.Expression()
				whereEXPRs = append(whereEXPRs, exp.Expression())
			}
		}
		ds = ds.Where(goqu.Or(whereEXPRs...))

		cnt := listCount
		if search.Options.Count.Valid {
			cnt = uint(search.Options.Count.Int64)
		}
		ds = ds.Limit(cnt)

		offset := uint(0)
		if search.Options.Page.Valid {
			offset = uint(search.Options.Page.Int64)
		}
		ds = ds.Offset(offset * cnt)
	}

	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	err = ds.ScanStructs(&bookResult)
	if err != nil {
		log.Println("ds: ", err.Error())
		return nil, err
	}
	if bookResult != nil {
		result = bookResult
	}

	return result, nil
}

// UpdateData - crUd
func UpdateData(data interface{}) (sql.Result, error) {
	whereEXP, err := CheckValidAndPrepareWhere(data)
	if err != nil {
		return nil, err
	}

	dbms := goqu.New("sqlite3", Dsn)
	ds := dbms.Update(TableName).Set(data).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dsn.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteData - cruD
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

// SelectCount - data count -> pages = (data count) / (count per page)
func SelectCount() (uint, error) {
	var cnt uint

	dbms := goqu.New("sqlite3", Dsn)
	ds := dbms.From(TableName).Select(goqu.COUNT("*").As("PAGE_COUNT"))
	ds.ScanVal(&cnt)

	return cnt, nil
}
