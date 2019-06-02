package schema

//go:generate env GOBIN=$PWD/.bin GO111MODULE=on go install github.com/sourcegraph/go-jsonschema/cmd/go-jsonschema-compiler
////go:generate $PWD/.bin/go-jsonschema-compiler -o schema.go -pkg schema aws_codecommit.schema.json bitbucket_server.schema.json critical.schema.json site.schema.json settings.schema.json github.schema.json  gitlab.schema.json gitolite.schema.json other_external_service.schema.json phabricator.schema.json
//go:generate $PWD/.bin/go-jsonschema-compiler -o schema.go -pkg schema aws_codecommit.schema.json

//go:generate env GO111MODULE=on go run stringdata.go -i aws_codecommit.schema.json -name AWSCodeCommitSchemaJSON -pkg schema -o aws_codecommit_stringdata.go
