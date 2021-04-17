module db

go 1.16

require (
	github.com/denisenkom/go-mssqldb v0.0.0-20200206145737-bbfc9a55622e
	github.com/doug-martin/goqu/v9 v9.11.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/lib/pq v1.2.0
	github.com/thoas/go-funk v0.8.0
	gopkg.in/guregu/null.v4 v4.0.0
	models v0.0.0
)

replace models => ../models
