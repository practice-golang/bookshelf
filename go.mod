module bookshelf

go 1.16

require (
	book v0.0.0
	consts v0.0.0
	db v0.0.0
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/labstack/echo/v4 v4.2.2
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/ini.v1 v1.62.0
	modernc.org/sqlite v1.10.1
)

replace (
	book => ./book
	consts => ./consts
	db => ./db
	models => ./models
)
