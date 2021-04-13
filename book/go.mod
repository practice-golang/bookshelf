module book

go 1.16

require (
	db v0.0.0
	github.com/labstack/echo/v4 v4.2.2
	models v0.0.0
)

replace (
	db => ../db
	models => ../models
)
