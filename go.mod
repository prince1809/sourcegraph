module github.com/prince1809/sourcegraph

go 1.12

require (
	github.com/fatih/color v1.7.0
	github.com/gchaincl/sqlhooks v1.1.0
	github.com/go-delve/delve v1.2.0 // indirect
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/google/zoekt v0.0.0-00010101000000-000000000000 // indirect
	github.com/lib/pq v1.1.1
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/mattn/goreman v0.0.0-00010101000000-000000000000 // indirect
	github.com/mcuadros/go-version v0.0.0-20190308113854-92cdf37c5b75
	github.com/opentracing-contrib/go-stdlib v0.0.0-20190324214902-3020fec0e66b
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.2
	github.com/sourcegraph/docsite v0.0.0-20190329030636-57dceb634057 // indirect
	github.com/sourcegraph/jsonx v0.0.0-20190114210550-ba8cb36a8614
	golang.org/x/net v0.0.0-20190110200230-915654e7eabc
	gopkg.in/inconshreveable/log15.v2 v2.0.0-20180818164646-67afb5ed74ec
)

replace (
	github.com/google/zoekt => github.com/sourcegraph/zoekt v0.0.0-20190116094554-c742ce874aa3
	github.com/graph-gophers/graphql-go => github.com/sourcegraph/graphql-go v0.0.0-20180929065141-c790ffc3c46a
	github.com/mattn/goreman => github.com/sourcegraph/goreman v0.1.2-0.20180928223752-6e9a2beb830d
	github.com/uber/gonduit => github.com/sourcegraph/gonduit v0.4.0
)

replace github.com/dghubble/gologin => github.com/sourcegraph/gologin v1.0.2-0.20181110030308-c6f1b62954d8

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
