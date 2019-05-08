package schema

// AWSCodeCommitConnection description: Configuration for a connection to AWS CodeCommit.
type AWSCodeCommitConnection struct {
	AccessKeyID                 string `json:"accessKeyID"`
	InitialRepositoryEnablement bool   `json:"initialRepositoryEnablement"`
	Region                      string `json:"region"`
	RepositoryPathPattern       string `json:"repositoryPathPattern"`
	SecretAccessKey             string `json:"secretAccessKey"`
}

// AuthAccessTokens description: Settings for access tokens, which enable external tools to access the Sourcegraph API with the privileges of the user.
type AuthAccessTokens struct {
	Allow string `json:"allow,omitempty"`
}

// AuthProviderCommon description: common properties for authentication providers.
type AuthProviderCommon struct {
	DisplayName string `json:"displayName,omitempty"`
}

// CriticalConfiguration description: Critical configuration for a Sourcegraph site.
type CriticalConfiguration struct {
	AuthDisableUsernameChanges bool `json:"auth.disableUsernameChanges,omitempty"`
	AuthEnableUsernameChanges  bool `json:"auth.enableUsernameChanges,omitempty"`
}

type OpenIDConnectAuthProvider struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	ConfigID     string `json:"configID,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	Issuer       string `json:"issuer"`
}

// OtherExternalServiceConnection description: Configuration for a Connection
type OtherExternalServiceConnection struct {
	Repos []string `json:"repos"`
	Url   string   `json:"url,omitempty"`
}

// ParentSourcegraph description: URL to fetch unreachable repository details from. Defaults to "https://sourcegraph.com"
type ParentSourcegraph struct {
}

// Phabricator description: Phabricator instance that integrates with this Gitolite instance.
type Phabricator struct {
}

// PhabricatorConnection description: Configuration for a connection to Phabricator.
type PhabricatorConnection struct {
}

type Repos struct {
}

// SAMLAuthProvider description: Configures the SAML authentication provider for SSO.
//
// Note: if you are using IdP-initiated login, you must have *at most one* SAMLAuthProvider in the `auth.providers` array.
type SAMLAuthProvider struct {
}

// SMTPServerConfig description: The SMTP server used to send transactional emails (such as email verification, reset-password emails, and notifications)
type SMTPServerConfig struct {
}

type SearchSavedQueries struct {
}

type SearchScope struct {
}

// Sentry description: Configuration for Sentry
type Sentry struct {
	Dsn string `json:"dsn,omitempty"`
}

// Settings description: Configuration settings for users and organization on Sourcegraph.
type Settings struct {
	Extensions map[string]bool `json:"extensions,omitempty"`
	Motd       []string        `json:"motd,omitempty"`
}

// SiteConfiguration description: Configuration for a Sourcegraph site.
type SiteConfiguration struct {
	AuthAccessTokens       *AuthAccessTokens `json:"auth.accessTokens,omitempty"`
	GitMaxConcurrentClones int               `json:"gitMaxConcurrentClones"`
	CorsOrigin             string            `json:"corsOrigin"`
}
