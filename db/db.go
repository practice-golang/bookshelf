package db // import "db"

import "database/sql"

var (
	Filename     = "./data.db"
	TableName    = "BOOKSHELF"
	Dsn          *sql.DB
	UpdateTarget []string // UPDATE ... WHERE idx=?
)

const (
	_ = iota
	// MYSQL - mysql
	MYSQL
	// SQLITE - sqlite
	SQLITE
)

// InitDB - DB파일 생성
func InitDB() (*sql.DB, error) {
	var err error
	fn := Filename

	Dsn, err = sql.Open("sqlite", fn)
	if err != nil {
		return nil, err
	}

	return Dsn, nil
}
