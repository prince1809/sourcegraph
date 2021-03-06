module github.com/prince1809/sourcegraph

go 1.12

require (
	github.com/Microsoft/go-winio v0.4.12 // indirect
	github.com/NYTimes/gziphandler v1.1.1
	github.com/alecthomas/chroma v0.6.3 // indirect
	github.com/certifi/gocertifi v0.0.0-20190506164543-d2eda7129713 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fatih/color v1.7.0
	github.com/gchaincl/sqlhooks v1.1.0
	github.com/getsentry/raven-go v0.2.0
	github.com/go-delve/delve v1.2.0
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/golang-migrate/migrate v3.5.4+incompatible
	github.com/golangci/golangci-lint v1.16.0
	github.com/google/uuid v1.1.1
	github.com/google/zoekt v0.0.0-00010101000000-000000000000
	github.com/gorilla/context v1.1.1
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.2
	github.com/hashicorp/go-multierror v1.0.0
	github.com/joho/godotenv v1.3.0
	github.com/keegancsmith/sqlf v1.1.0
	github.com/keegancsmith/tmpfriend v0.0.0-20180423180255-86e88902a513
	github.com/kevinburke/differ v0.0.0-20180721181420-bdfd927653c8
	github.com/kevinburke/go-bindata v3.13.0+incompatible
	github.com/kr/text v0.1.0
	github.com/lib/pq v1.1.1
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/mattn/goreman v0.0.0-00010101000000-000000000000
	github.com/mcuadros/go-version v0.0.0-20190308113854-92cdf37c5b75
	github.com/neelance/parallel v0.0.0-20160708114440-4de9ce63d14c
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
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
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0
	golang.org/x/net v0.0.0-20190313220215-9f648a60d977
	golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4
	golang.org/x/sys v0.0.0-20190502145724-3ef323f4f1fd
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
