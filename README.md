# Practice
```
golang 1.16 embed, doug-martin/goqu, cznic/sqlite CRUD
```

## Build & Run
* `make` -> `cd bin` -> `bookshelf` or `go build` -> `bookshelf` or `go run main.go`
* Read/Write test : `requests.http`

## Preference file
* Once on `bookshelf` binary first run, `bookshelf.ini` which contain db connection will be created

## Database
* Default set is `sqlite`
* See `consts/db.go` or generated `bookshelf.ini`
  * `sqlite`, `mysql`, `sqlserver`, `postrges`

## Init sample data & Open the webpage
* Find and send request `### Add book #1` in `requests.http`
* Open http://localhost:2918

## Embeded html
* `index.html` - Used `vue.js`, `sakura.css`

## Test and Todo
* No. YOLO!😆


## License
Public domain
