module github.com/prince1809/sourcegraph

go 1.12

require (
	github.com/alecthomas/chroma v0.6.3 // indirect
	github.com/fatih/color v1.7.0
	github.com/gchaincl/sqlhooks v1.1.0
	github.com/go-delve/delve v1.2.0
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/golangci/golangci-lint v1.16.0
	github.com/google/zoekt v0.0.0-00010101000000-000000000000
	github.com/hashicorp/go-multierror v1.0.0
	github.com/keegancsmith/sqlf v1.1.0
	github.com/kevinburke/differ v0.0.0-20180721181420-bdfd927653c8
	github.com/kevinburke/go-bindata v3.13.0+incompatible
	github.com/labstack/gommon v0.2.8
	github.com/lib/pq v1.1.1
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/mattn/goreman v0.0.0-00010101000000-000000000000
	github.com/mcuadros/go-version v0.0.0-20190308113854-92cdf37c5b75
	github.com/opentracing-contrib/go-stdlib v0.0.0-20190324214902-3020fec0e66b
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.3
	github.com/prometheus/common v0.4.0
	github.com/shurcooL/httpfs v0.0.0-20181222201310-74dc9339e414
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sourcegraph/docsite v0.0.0-20190329030636-57dceb634057
	github.com/sourcegraph/go-jsonschema v0.0.0-20190205151546-7939fa138765
	github.com/sourcegraph/jsonx v0.0.0-20190114210550-ba8cb36a8614
	github.com/valyala/fasttemplate v1.0.1 // indirect
	golang.org/x/net v0.0.0-20190313220215-9f648a60d977
	golang.org/x/sys v0.0.0-20190502145724-3ef323f4f1fd // indirect
	golang.org/x/tools v0.0.0-20190314010720-f0bfdbff1f9c
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
