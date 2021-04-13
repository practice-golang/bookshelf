package main // import "bookshelf"

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	"time"

	_ "modernc.org/sqlite"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"book"
	"db"
)

var (
	//go:embed static
	content embed.FS
)

func setupDB() error {
	var err error

	db.Dsn, err = db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	recreate := false
	err = db.CreateTable(recreate)
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
	e.GET("/books", book.GetBooks)
	e.PUT("/books", book.AddBooks)
	e.PATCH("/books", book.EditBooks)
	e.DELETE("/book/:idx", book.DeleteBook)

	return e
}

func main() {
	var fileConnectionLog *os.File
	var err error

	db.UpdateTarget = []string{"idx"} // UPDATE ... WHERE idx=?

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
