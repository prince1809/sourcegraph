package api

type RepoID int32

type RepoName string

type CommitID string

type Repo struct {
	ID           RepoID
	ExternalRepo *ExternalRepoSpec
	Name         RepoName
	Enabled      bool
}

type InsertRepoOp struct {
	Name         RepoName
	Description  string
	Fork         bool
	Archived     bool
	Enabled      bool
	ExternalRepo *ExternalRepoSpec
}

type ExternalRepoSpec struct {
	ID          string
	ServiceType string
	ServiceID   string
}
