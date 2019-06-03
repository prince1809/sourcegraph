package db

import (
	"errors"
	"time"
)

// AccessToken describes an access token. The actual token (that a caller must supply to
// authenticate) is not stored and is not present in this struct.
type AccessToken struct {
	ID            int64
	SubjectUserID int32 // the user whose privileges the access token grants.
	Scopes        []string
	Note          string
	CreatorUserID int32
	CreatedAt     time.Time
	LastUsedAt    *time.Time
}

// ErrAccessTokenNotFound occurs when a database operation expects a specific access token to exist
// but it does not exists.
var ErrAccessTokenNotFound = errors.New("access token not found")

// accessTokens implements autocert.Cache
type accessTokens struct{}


// Create creates an access token for the specified user. The secret token value itself is
// returned. The caller is responsible for presenting this value to the end user; Sourcegraph does
// not retain it (only a hash of it).
//
// The secret token value is a long random string; it is what API client must provide to
// authenticate their requests.
