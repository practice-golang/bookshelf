module db

go 1.16

require (
	github.com/doug-martin/goqu/v9 v9.11.0
	github.com/google/go-cmp v0.5.5
	github.com/thoas/go-funk v0.8.0
	gopkg.in/guregu/null.v4 v4.0.0
	models v0.0.0
)

replace models => ../models
