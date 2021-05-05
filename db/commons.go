package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/practice-golang/bookshelf/models"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
)

func getDialect() (string, error) {
	var dbType string

	switch DBType {
	case SQLITE:
		dbType = "sqlite3"
	case MYSQL:
		dbType = "mysql"
	case POSTGRES:
		dbType = "postgres"
	case SQLSERVER:
		dbType = "sqlserver"
	default:
		return dbType, errors.New("nothing to support DB")
	}

	return dbType, nil
}

// InsertData - Crud
func InsertData(data interface{}) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Insert(TableName).Rows(data)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// InsertDataMAP - Crud
func InsertDataMAP(data interface{}) (sql.Result, error) {
	rcds := []goqu.Record{}
	if data != nil {
		jsonBody, ok := gjson.Parse(string(data.([]byte))).Value().([]interface{})
		if !ok {
			log.Println("Cannot parse jsonBody")
		}

		for i, d := range jsonBody {
			log.Println(i, d)
			rcd := goqu.Record{}
			for j, v := range d.(map[string]interface{}) {
				rcd[j] = v
			}
			rcds = append(rcds, rcd)
		}

	}

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Insert(TableName).Rows(rcds)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SelectData - cRud
func SelectData(search interface{}) (interface{}, error) {
	var result interface{}
	var err error

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.From(TableName).Select(search)

	bookResult := []models.Book{}

	switch search := search.(type) {
	case models.Book:
		exps := PrepareWhere(search)
		if !exps.IsEmpty() {
			ds = ds.Where(exps)
		}
		// cnt := listCount
		cnt := uint(10)
		ds = ds.Limit(cnt)

	case models.BookSearch:
		keywords := search.Keywords
		exps := []goqu.Expression{}
		for _, k := range keywords {
			ex := PrepareWhere(k)
			if !ex.IsEmpty() {
				for c, v := range ex {
					val := fmt.Sprintf("%s%s%s", "%", v, "%")
					ex[c] = goqu.Op{"like": val}
				}
				exps = append(exps, ex.Expression())
			}
		}
		ds = ds.Where(goqu.Or(exps...))

		orderDirection := goqu.C(OrderScope).Asc()
		if search.Options.Order.String == "desc" {
			orderDirection = goqu.C(OrderScope).Desc()
		}
		ds = ds.Order(orderDirection)

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

// SelectDataMAP - cRud MAP
func SelectDataMAP(search interface{}) (interface{}, error) {
	var result []map[string]interface{}
	var err error

	offset := uint(0)
	count := listCount

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	colNames := []interface{}{"NAME", "AUTHOR"}
	colNames = append(colNames, "IDX")
	colNames = append(colNames, "PRICE")
	colNames = append(colNames, "ISBN")
	dbms := goqu.New(dbType, Dbo)
	ds := dbms.From(TableName).Select(colNames...)

	exps := []goqu.Expression{}
	if search != nil {
		jsonBody, ok := gjson.Parse(string(search.([]byte))).Value().(map[string]interface{})
		if !ok {
			log.Println("Cannot parse jsonBody")
		}

		if jsonBody["keywords"] != nil {
			for i, kwds := range jsonBody["keywords"].([]interface{}) {
				log.Println("Keywords: ", i, kwds)
				ex := goqu.Ex{}

				for k, v := range kwds.(map[string]interface{}) {
					ex[k] = goqu.Op{"like": v}
				}
				if !ex.IsEmpty() {
					exps = append(exps, ex.Expression())
				}
			}
		}

		if jsonBody["options"] != nil {
			if jsonBody["options"].(map[string]interface{})["count"] != nil {
				count = uint(jsonBody["options"].(map[string]interface{})["count"].(float64))
			}

			if jsonBody["options"].(map[string]interface{})["page"] != nil {
				offset = uint(jsonBody["options"].(map[string]interface{})["page"].(float64))
			}
		}
	}

	ds = ds.Where(goqu.Or(exps...))

	cnt := uint(count)
	ds = ds.Limit(cnt)

	ds = ds.Offset(offset * cnt)

	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	rows, err := Dbo.Query(sql, args...)
	if err != nil {
		log.Println("rowsMAP: ", err.Error())
		return nil, err
	}
	cols, _ := rows.Columns()

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		colPtrs := make([]interface{}, len(cols))
		for i := range columns {
			colPtrs[i] = &columns[i]
		}

		if err := rows.Scan(colPtrs...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := colPtrs[i].(*interface{})
			switch (*val).(type) {
			case string:
				m[colName] = (*val).(string)
			case uint8:
				m[colName] = (*val).(uint8)
			case []uint8:
				// string or double or integer or numeric
				m[colName] = string([]byte((*val).([]uint8)))
			case int64:
				m[colName] = (*val).(int64)
			case float64:
				m[colName] = (*val).(float64)
			}
		}

		result = append(result, m)
	}

	return result, err
}

// UpdateData - crUd
func UpdateData(data interface{}) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	whereEXP, err := CheckValidAndPrepareWhere(data)
	if err != nil {
		return nil, err
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Update(TableName).Set(data).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateDataMAP - crUd
func UpdateDataMAP(data interface{}) (sql.Result, error) {
	whereEXP := goqu.Ex{}
	record := goqu.Record{}
	if data != nil {
		jsonBody, ok := gjson.Parse(string(data.([]byte))).Value().(map[string]interface{})
		if !ok {
			log.Println("Cannot parse jsonBody")
		}

		for k, v := range jsonBody {
			record[k] = v

			if funk.Contains(UpdateScope, k) {
				whereEXP[k] = goqu.Op{"eq": v}
			}
		}
	}

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Update(TableName).Set(record).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteData - cruD
func DeleteData(target, value string) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	whereEXP := goqu.Ex{target: value}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Delete(TableName).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SelectCount - data count -> pages = (data count) / (count per page)
func SelectCount(search interface{}) (uint, error) {
	var cnt uint

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.From(TableName).Select(goqu.COUNT("*").As("PAGE_COUNT"))

	switch search := search.(type) {
	case models.BookSearch:
		keywords := search.Keywords
		exps := []goqu.Expression{}
		for _, k := range keywords {
			ex := PrepareWhere(k)
			if !ex.IsEmpty() {
				for c, v := range ex {
					val := fmt.Sprintf("%s%s%s", "%", v, "%")
					ex[c] = goqu.Op{"like": val}
				}
				exps = append(exps, ex.Expression())
			}
		}
		ds = ds.Where(goqu.Or(exps...))
	}

	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	ds.ScanVal(&cnt)

	return cnt, nil
}
