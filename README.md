# Practice
```
golang embed, doug-martin/goqu, cznic/sqlite CRUD
```

## Build & Run
* `make` -> `cd bin` -> `bookshelf` or `go build` -> `bookshelf` or `go run main.go`
* Read/Write test : `requests.http`

## Database
* Default is `sqlite`
* See `consts/db.go` to choose one
  * `sqlite`, `mysql`, `sqlserver`, `postrges`

## Start
* Find and send request `### Add book #1` in `requests.http`
* Open http://localhost:2918
