package main // import "bookshelf"

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Book struct {
	IDX       int64   `json:"idx"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	SalePrice float64 `json:"sale-price"`
	PriceUnit string  `json:"price-unit"`
	Author    string  `json:"author"`
	Edition   int64   `json:"edition"`
	Publisher string  `json:"publisher"`
	ISBN      string  `json:"isbn"`
}

var (
	dbFileName = "./data.db"
	tableName  = "BOOKSHELF"
	db         *sql.DB
	//go:embed static/*
	content embed.FS
)

var booksDummy = []Book{
	{
		Name:      "흔한남매 7",
		Price:     10800,
		SalePrice: 9300,
		PriceUnit: "KRW",
		Author:    "백난도",
		Edition:   1,
		Publisher: "미래엔아이세움",
		ISBN:      "9791164137527",
	},
	{
		Name:      "성장의 종말",
		Price:     17000,
		SalePrice: 15300,
		PriceUnit: "KRW",
		Author:    "디트리히 볼래스",
		Edition:   1,
		Publisher: "더퀘스트",
		ISBN:      "9791165215170",
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
		"SALE_PRICE"	REAL,
		"PRICE_UNIT"	TEXT,
		"AUTHOR"		TEXT,
		"EDITION"		INTEGER,
		"PUBLISHER"		TEXT,
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
	INSERT OR REPLACE INTO "#TABLE_NAME"
		(NAME, PRICE, SALE_PRICE, PRICE_UNIT, AUTHOR, PUBLISHER, EDITION, ISBN)
	VAlUES("#BOOK_NAME", #PRICE_NORMAL, #SALE_PRICE, "#PRICE_UNIT", "#AUTHOR", "#PUBLISHER", #EDITION, #ISBN);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableName)

	sql = strings.ReplaceAll(sql, "#BOOK_NAME", book.Name)
	sql = strings.ReplaceAll(sql, "#PRICE_NORMAL", fmt.Sprint(book.Price))
	sql = strings.ReplaceAll(sql, "#SALE_PRICE", fmt.Sprint(book.SalePrice))
	sql = strings.ReplaceAll(sql, "#PRICE_UNIT", book.PriceUnit)
	sql = strings.ReplaceAll(sql, "#AUTHOR", book.Author)
	sql = strings.ReplaceAll(sql, "#PUBLISHER", book.Publisher)
	sql = strings.ReplaceAll(sql, "#EDITION", fmt.Sprint(book.Edition))
	sql = strings.ReplaceAll(sql, "#ISBN", book.ISBN)

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// SelectData - cRud
func SelectData(db *sql.DB, search Book) ([]Book, error) {
	sql := ""

	sql += `
	SELECT
		IDX, NAME, PRICE, SALE_PRICE, PRICE_UNIT, AUTHOR, EDITION, PUBLISHER, ISBN
	FROM #TABLE_NAME
	`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableName)

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	result := []Book{}

	for rows.Next() {
		var b Book

		err = rows.Scan(&b.IDX, &b.Name, &b.Price, &b.SalePrice, &b.PriceUnit, &b.Author, &b.Edition, &b.Publisher, &b.ISBN)
		if err != nil {
			return nil, err
		}

		result = append(result, b)
	}

	return result, nil
}

// GetBooks - 책정보 취득
func GetBooks(c echo.Context) error {
	data, err := SelectData(db, Book{})
	if err != nil {
		log.Fatal("SelectData: ", err)
	}

	// for _, b := range data {
	// 	log.Println(b.IDX, b.Name, b.Author, b.Price)
	// }

	return c.JSON(http.StatusOK, data)
}

func serverSetup() *echo.Echo {
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

	return e
}

func main() {
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

	for _, book := range booksDummy {
		err = InsertData(book)
		if err != nil {
			log.Fatal("InsertData: ", err)
		}
	}

	e := serverSetup()
	e.Logger.Fatal(e.Start("127.0.0.1:2918"))
}
