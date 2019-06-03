package db

import "github.com/prince1809/sourcegraph/schema"

// An ExternalServicesStore stores external services and their configuration.
// Before updating or creating a new external service, validation is performed.
// The enterprise code registers additional validators at run-time and sets the
// global instance in stores.go
type ExternalServicesStore struct {
	GithubValidators []func(*schema.GitHubConnection) error
	GitLabValidators []func(*schema.GitLabConnection, []schema.AuthProviders)
}

// ExternalServicesListOptions contains options for listing external services.
type ExternalServicesListOptions struct {
	Kinds []string
	*LimitOffset
}
