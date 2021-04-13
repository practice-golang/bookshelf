package main // import "bookshelf"

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "modernc.org/sqlite"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gopkg.in/guregu/null.v4"

	"github.com/thoas/go-funk"
)

type Book struct {
	Idx    null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	Name   null.String `json:"name" db:"NAME"`
	Price  null.Float  `json:"price" db:"PRICE"`
	Author null.String `json:"author" db:"AUTHOR"`
	ISBN   null.String `json:"isbn" db:"ISBN"`
}

var (
	dbFileName = "./data.db"
	tableName  = "BOOKSHELF"
	db         *sql.DB
	//go:embed static
	content      embed.FS
	updateTarget = []string{"idx"} // UPDATE ... WHERE idx=?
)

// InitDB - DB파일 생성
func InitDB() (*sql.DB, error) {
	fn := dbFileName

	db, err := sql.Open("sqlite", fn)
	if err != nil {
		return nil, err
	}

	return db, nil
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

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableName)

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CheckValidAndPrepareWhere - 빠진 요소 체크, Where 요소 준비
func CheckValidAndPrepareWhere(book Book) (goqu.Ex, error) {
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
			if funk.Contains(updateTarget, jsonName) {
				result[jsonName], _ = f.Value()
			}
		case null.Int:
			if !f.Valid {
				return nil, errors.New("`" + jsonName + "` must have a value")
			}
			if funk.Contains(updateTarget, jsonName) {
				result[jsonName], _ = f.Value()
			}
		case null.Float:
			if !f.Valid {
				return nil, errors.New("`" + jsonName + "` must have a value")
			}
			if funk.Contains(updateTarget, jsonName) {
				result[jsonName], _ = f.Value()
			}
		}
	}

	return result, nil
}

// InsertData - Crud
func InsertData(books []Book) (sql.Result, error) {
	dbms := goqu.New("sqlite3", db)
	ds := dbms.Insert(tableName).Rows(books)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := db.Exec(sql)
	if err != nil {
		// lastID, _ := result.LastInsertId()
		// affRows, _ := result.RowsAffected()
		// log.Println(lastID, affRows)
		return nil, err
	}

	return result, nil
}

// SelectData - cRud
func SelectData(search Book) ([]Book, error) {
	var result []Book

	dbms := goqu.New("sqlite3", db)
	ds := dbms.From(tableName).Select(&Book{})
	err := ds.ScanStructs(&result)
	if err != nil {
		log.Println("ds: ", err.Error())
		return nil, err
	}

	return result, nil
}

// UpdateData - crUd
func UpdateData(book Book) (sql.Result, error) {
	whereEXP, err := CheckValidAndPrepareWhere(book)
	if err != nil {
		return nil, err
	}

	dbms := goqu.New("sqlite3", db)
	ds := dbms.Update(tableName).Set(book).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := db.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteData(target, value string) (sql.Result, error) {
	whereEXP := goqu.Ex{target: value}

	dbms := goqu.New("sqlite3", db)
	ds := dbms.Delete(tableName).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := db.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// AddBooks - 책정보 입력
func AddBooks(c echo.Context) error {
	var books []Book

	if err := c.Bind(&books); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	sqlResult, err := InsertData(books)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// GetBooks - 책정보 취득
func GetBooks(c echo.Context) error {
	data, err := SelectData(Book{})
	if err != nil {
		log.Fatal("SelectData: ", err)
	}

	// for _, b := range data {
	// 	log.Println(b.IDX, b.Name, b.Author, b.Price)
	// }

	return c.JSON(http.StatusOK, data)
}

// EditBooks - 책정보 수정
func EditBooks(c echo.Context) error {
	var book Book

	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	sqlResult, err := UpdateData(book)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteBook - 책 1개 삭제
func DeleteBook(c echo.Context) error {
	idx := c.Param("idx")

	sqlResult, err := DeleteData("idx", idx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

func setupDB() error {
	var err error

	db, err = InitDB()
	if err != nil {
		log.Fatal(err)
	}

	recreate := false
	err = CreateTable(recreate)
	if err != nil {
		log.Fatal("CreateTable: ", err)
	}

	return err
}

func dumpHandler(c echo.Context, reqBody, resBody []byte) {
	header := time.Now().Format("2006-01-02 15:04:05") + " - "
	body := string(reqBody)
	body = strings.Replace(body, "\r\n", "", -1)
	body = strings.Replace(body, "\n", "", -1)
	data := header + body + "\n"

	f, err := os.OpenFile(
		"request-body.log",
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		os.FileMode(0777),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		log.Println(err)
		return
	}
}

func setupServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(
		middleware.CORS(),
		middleware.Recover(),
	)

	contentHandler := echo.WrapHandler(http.FileServer(http.FS(content)))
	contentRewrite := middleware.Rewrite(map[string]string{"/*": "/static/$1"})

	e.GET("/*", contentHandler, contentRewrite)
	e.GET("/books", GetBooks)
	e.PUT("/books", AddBooks)
	e.PATCH("/books", EditBooks)
	e.DELETE("/book/:idx", DeleteBook)

	return e
}

func main() {
	var fileConnectionLog *os.File
	var err error

	err = setupDB()
	if err != nil {
		log.Fatal("Setup DB: ", err)
	}

	e := setupServer()

	fileConnectionLog, err = os.OpenFile(
		"connection.log",
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		os.FileMode(0777),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer fileConnectionLog.Close()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339} - remote_ip:${remote_ip}, host:${host}, ` +
			`method:${method}, uri:${uri},status:${status}, error:${error}, ` +
			`${header:Authorization}, query:${query:property}, form:${form}, ` + "\n",
		Output: fileConnectionLog,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		AllowMethods: []string{
			echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE,
			echo.HEAD, echo.OPTIONS,
		},
	}))

	e.Use(middleware.BodyDump(dumpHandler))

	e.Logger.Fatal(e.Start("127.0.0.1:2918"))
}
