package main // import "bookshelf"

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "modernc.org/sqlite"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gopkg.in/guregu/null.v4"
)

type Book struct {
	Idx    null.String `json:"idx" db:"IDX"`
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
	content embed.FS
)

var booksDummy = []Book{
	{
		Name:   null.NewString("흔한남매 7", true),
		Price:  null.NewFloat(10800, true),
		Author: null.NewString("백난도", true),
		ISBN:   null.NewString("9791164137527", true),
	},
	{
		Name:   null.NewString("성장의 종말", true),
		Price:  null.NewFloat(17000, true),
		Author: null.NewString("디트리히 볼래스", true),
		ISBN:   null.NewString("9791165215170", true),
	},
}

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

// InsertData - Crud
func InsertData(book Book) error {
	sql := ""

	sql += `
	INSERT INTO "#TABLE_NAME"
		(NAME, PRICE, AUTHOR, ISBN)
	VAlUES("#BOOK_NAME", #PRICE_NORMAL, "#AUTHOR", #ISBN);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableName)

	sql = strings.ReplaceAll(sql, "#BOOK_NAME", book.Name.String)
	sql = strings.ReplaceAll(sql, "#PRICE_NORMAL", fmt.Sprint(book.Price.Float64))
	sql = strings.ReplaceAll(sql, "#AUTHOR", book.Author.String)
	sql = strings.ReplaceAll(sql, "#ISBN", book.ISBN.String)

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// InsertReplaceData - Crud
func InsertReplaceData(book Book) error {
	sql := ""

	sql += `
	INSERT OR REPLACE INTO "#TABLE_NAME"
		(NAME, PRICE, AUTHOR, ISBN)
	VAlUES("#BOOK_NAME", #PRICE_NORMAL, "#AUTHOR", #ISBN);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableName)

	sql = strings.ReplaceAll(sql, "#BOOK_NAME", book.Name.String)
	sql = strings.ReplaceAll(sql, "#PRICE_NORMAL", fmt.Sprint(book.Price.Float64))
	sql = strings.ReplaceAll(sql, "#AUTHOR", book.Author.String)
	sql = strings.ReplaceAll(sql, "#ISBN", book.ISBN.String)

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
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

// SetDummy - 더미데이터 입력
func SetDummy(c echo.Context) error {
	var err error
	for _, book := range booksDummy {
		err = InsertData(book)
		if err != nil {
			log.Println("InsertData: ", err)
			return c.JSON(http.StatusConflict, map[string]string{"msg": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, "OK")
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
	e.PUT("/books/dummy", SetDummy)

	return e
}

func main() {
	err := setupDB()
	if err != nil {
		log.Fatal("Setup DB: ", err)
	}

	e := setupServer()
	e.Logger.Fatal(e.Start("127.0.0.1:2918"))
}
