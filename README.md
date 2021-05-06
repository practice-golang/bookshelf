# Practice
```
golang 1.16 embed, doug-martin/goqu, cznic/sqlite CRUD
```

## Build & Run
* `go get github.com/practice-golang/bookshelf` -> run `bookshelf` or
* `make` -> `cd bin` -> run `bookshelf`
* Read/Write test : `requests.http`, `requests-map.http`

## Preference file
* Once on `bookshelf` binary first run, `bookshelf.ini` which contain db connection will be created

## Database
* Default set in `bookshelf.ini` is `sqlite`
* See `config/db.go` or generated `bookshelf.ini`
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
