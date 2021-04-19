package main // import "github.com/practice-golang/bookshelf"

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"time"

	_ "modernc.org/sqlite"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gopkg.in/ini.v1"

	"github.com/practice-golang/bookshelf/book"
	"github.com/practice-golang/bookshelf/consts"
	"github.com/practice-golang/bookshelf/db"
)

var (
	//go:embed static
	content embed.FS
	//go:embed samples/bookshelf.ini
	sampleINI string
)

func setupDB() error {
	var err error

	switch consts.DbInfo.Type {
	case "sqlite":
		db.DBType = db.SQLITE
		db.Dsn = consts.DbInfo.Filename
	case "mysql":
		db.DBType = db.MYSQL
		db.Dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/",
			consts.DbInfo.User,
			consts.DbInfo.Password,
			consts.DbInfo.Server,
			consts.DbInfo.Port,
		)
		db.DatabaseName = consts.DbInfo.Database
		db.TableName = db.DatabaseName + "." + db.TableName
	case "postgres":
		db.DBType = db.POSTGRES
		db.Dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", consts.DbInfo.Server, consts.DbInfo.Port, consts.DbInfo.User, consts.DbInfo.Password, consts.DbInfo.Database)

		db.DatabaseName = consts.DbInfo.Schema
		db.TableName = db.DatabaseName + "." + db.TableName
	case "sqlserver":
		db.DBType = db.SQLSERVER
		db.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", consts.DbInfo.User, consts.DbInfo.Password, consts.DbInfo.Server, consts.DbInfo.Port, consts.DbInfo.Database)
		db.DatabaseName = consts.DbInfo.Database
		db.TableName = db.DatabaseName + ".dbo." + db.TableName
	default:
		log.Fatal("nothing to support DB")
	}

	db.Dbi, err = db.InitDB(db.DBType)
	if err != nil {
		log.Fatal("InitDB: ", err)
	}

	recreate := false
	err = db.Dbi.CreateTable(recreate)
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
	e.GET("/book/:idx", book.GetBook)
	e.GET("/books", book.GetBooks)
	e.POST("/books", book.SearchBooks)
	e.PUT("/books", book.AddBooks)
	e.PATCH("/book", book.EditBook)
	e.DELETE("/book/:idx", book.DeleteBook)

	e.POST("/total-page", book.GetTotalPage)

	return e
}

func main() {
	var err error
	cfg, err := ini.Load("bookshelf.ini")

	if err != nil {
		log.Print("Fail to read ini. ")

		f, err := os.Create("bookshelf.ini")
		if err != nil {
			log.Fatal("Create INI: ", err)
		}
		defer f.Close()

		_, err = f.WriteString(sampleINI + "\n")
		if err != nil {
			log.Fatal("Create INI: ", err)
		}

		log.Println("bookshelf.ini is created")
	}

	if cfg != nil {
		consts.DbInfo.Type = cfg.Section("database").Key("DBTYPE").String()
		consts.DbInfo.Server = cfg.Section("database").Key("ADDRESS").String()
		consts.DbInfo.Port, _ = cfg.Section("database").Key("PORT").Int()
		consts.DbInfo.User = cfg.Section("database").Key("USER").String()
		consts.DbInfo.Password = cfg.Section("database").Key("PASSWORD").String()
		consts.DbInfo.Database = cfg.Section("database").Key("DATABASE").String()
		consts.DbInfo.Schema = cfg.Section("database").Key("SCHEMA").String()
		consts.DbInfo.Filename = cfg.Section("database").Key("FILENAME").String()
	}

	var fileConnectionLog *os.File

	db.UpdateScope = []string{"IDX"}             // UPDATE ... WHERE IDX=?
	db.IgnoreScope = []string{"AUTHOR", "PRICE"} // Ignore if nil or null
	db.OrderScope = "IDX"

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

	// e.Logger.Fatal(e.Start(":2918"))
	e.Logger.Fatal(e.Start("127.0.0.1:2918"))
}
