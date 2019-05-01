module github.com/prince1809/sourcegraph

go 1.12

require (
	github.com/go-delve/delve v1.2.0 // indirect
	github.com/google/zoekt v0.0.0-20190425223102-3493be98f11f // indirect
	github.com/mattn/goreman v0.2.1 // indirect
	github.com/mcuadros/go-version v0.0.0-20190308113854-92cdf37c5b75
	github.com/sourcegraph/docsite v0.0.0-20190329030636-57dceb634057 // indirect
)

replace (
	github.com/google/zoekt => github.com/sourcegraph/zoekt v0.0.0-20190116094554-c742ce874aa3
	github.com/graph-gophers/graphql-go => github.com/sourcegraph/graphql-go v0.0.0-20180929065141-c790ffc3c46a
	github.com/mattn/goreman => github.com/sourcegraph/goreman v0.1.2-0.20180928223752-6e9a2beb830d
	github.com/uber/gonduit => github.com/sourcegraph/gonduit v0.4.0
)

replace github.com/dghubble/gologin => github.com/sourcegraph/gologin v1.0.2-0.20181110030308-c6f1b62954d8

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
