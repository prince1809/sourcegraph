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

type AuthProviders struct {
}

type BitbucketServerConnection struct {
}

type BrandAssets struct {
}

type Branding struct {
}

type BuilinAuthProvider struct {
}

type CloneURLToRepositoryName struct {
}

// CriticalConfiguration description: Critical configuration for a Sourcegraph site.
type CriticalConfiguration struct {
	AuthDisableUsernameChanges bool `json:"auth.disableUsernameChanges,omitempty"`
	AuthEnableUsernameChanges  bool `json:"auth.enableUsernameChanges,omitempty"`
}

// Discussions description: Configures Sourcegraph code discussions.
type Discussions struct {
	AbuseEmails     []string `json:"abuseEmails,omitempty"`
	AbuseProtection bool     `json:"abuseProtection,omitempty"`
}

type ExcludedBitbucketServerRepo struct {
}

type ExcludedGitHubRepo struct {
}

type ExcludedGitLabProject struct {
}

// ExperimentalFeatures description: Experimental features to enable or disable. Features that are now enabled by default are marked as deprecated.
type ExperimentalFeatures struct {
	Discussions string `json:"discussions,omitempty"`
}

type Extensions struct {
}

type ExternalIdentity struct {
}

type GithubAuthProvider struct {
}

type GithubAuthorization struct {
}

type GithubConnection struct {
}

type GitlabProvider struct {
}

type GitlabAuthorization struct {
}

type GitlabConnection struct {
}

type GitlabProject struct {
}

type GitoliteConnection struct {
}

type HTTPHeaderAuthProvider struct {
}

type IMAPServerConfig struct {
}

type IdentityProvider struct {
}

// Log description: Configuration for logging and alerting, including to external services.
type Log struct {
	Sentry *Sentry `json:"sentry,omitempty"`
}

type Notice struct {
	Dismissible bool   `json:"dismissible,omitempty"`
	Location    string `json:"location"`
	Message     string `json:"message"`
}

type OAuthIdentity struct {
	Type string `json:"type"`
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
	Authentication string `json:"authentication"`
	Domain         string `json:"domain,omitempty"`
	Host           string `json:"host"`
	Password       string `json:"password,omitempty"`
	Port           int    `json:"port"`
	Username       string `json:"username,omitempty"`
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
	AuthAccessTokens                  *AuthAccessTokens           `json:"auth.accessTokens,omitempty"`
	Branding                          *Branding                   `json:"branding,omitempty"`
	CorsOrigin                        string                      `json:"corsOrigin"`
	DisableAutoGitUpdates             bool                        `json:"disableAutoGitUpdates,omitempty"`
	DisableBuiltInSearches            bool                        `json:"disableBuiltInSearches"`
	DisablePublicRepoRedirects        bool                        `json:"disablePublicRepoRedirects,omitempty"`
	Discussions                       *Discussions                `json:"discussions,omitempty"`
	DontIncludeSymbolResultsByDefault bool                        `json:"dontIncludeSymbolResultsByDefault,omitempty"`
	EmailAddress                      string                      `json:"email.address,omitempty"`
	EmailImap                         *IMAPServerConfig           `json:"email.imap,omitempty"`
	EmailSmtp                         *SMTPServerConfig           `json:"email.smtp,omitempty"`
	ExperimentalFeatures              *ExperimentalFeatures       `json:"experimentalFeatures,omitempty"`
	Extensions                        *Extensions                 `json:"extensions,omitempty"`
	GitCloneURLToRepositoryName       []*CloneURLToRepositoryName `json:"git.cloneURLToRepositoryName,omitempty"`
	GitMaxConcurrentClones            int                         `json:"gitMaxConcurrentClones"`
	GithubClientID                    string                      `json:"githubClientID,omitempty"`
	GithubClientSecret                string                      `json:"githubClientSecret,omitempty"`
	MaxReposToSearch                  int                         `json:"maxReposToSearch,omitempty"`
	ParentSourcegraph                 *ParentSourcegraph          `json:"parentSourcegraph,omitempty"`
	RepoListUpdateInterval            int                         `json:"repoListUpdateInterval,omitempty"`
	SearchIndexEnabled                *bool                       `json:"search.index.enabled,omitempty"`
	SearchLargeFiles                  []string                    `json:"search.largeFiles,omitempty"`
}

// SlackNotificationsConfig description: Configuration for sending notification to slack.
type SlackNotificationConfig struct {
	WebhookURL string `json:"webhookURL"`
}

type UsernameIdentity struct {
	Type string `json:"type"`
}
