package consts // import "consts"

type DBpath struct {
	Type     string
	Server   string // mysql, postgres, sqlserver
	Port     int    // mysql, postgres, sqlserver
	User     string // mysql, postgres, sqlserver
	Password string // mysql, postgres, sqlserver
	Database string // mysql, postgres, sqlserver
	Schema   string // postgres
	Filename string // sqlite
}

// DB - sqlite
var (
	DbInfo = DBpath{
		Type:     "sqlite", // "sqlite" "mysql" "postgres" "sqlserver"
		Filename: "./bookshelf.db",
	}
)

// DB - mysql
// var (
// 	DbInfo = DBpath{
// 		Type:     "mysql", // "sqlite" "mysql" "postgres" "sqlserver"
// 		Server:   "127.0.0.1",
// 		Port:     13306,
// 		User:     "root",
// 		Password: "",
// 		Database: "bookshelf",
// 		Schema:   "",
// 		Filename: "./bookshelf.db",
// 	}
// )

// var (
// 	DbInfo = DBpath{
// 		Type:     "sqlserver", // "sqlite" "mysql" "postgres" "sqlserver"
// 		Server:   "127.0.0.1",
// 		Port:     1433,
// 		User:     "sa",
// 		Password: "mssql",
// 		Database: "bookshelf",
// 		Schema:   "",
// 	}
// )