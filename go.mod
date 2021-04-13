module bookshelf

go 1.16

require (
	book v0.0.0
	db v0.0.0
	github.com/labstack/echo/v4 v4.2.2
	modernc.org/sqlite v1.10.1
)

replace (
	book => ./book
	db => ./db
	models => ./models
)
