package gitserver

import "github.com/prince1809/sourcegraph/pkg/api"

type Cmd struct {
	//client *Client
}

type Repo struct {
	Name api.RepoName
	URL  string
}
