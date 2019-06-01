// Package actor provides the structures for representing an actor who has
// to access to resources.
package actor

import (
	"context"
	"github.com/prince1809/sourcegraph/pkg/trace"
)

// Actor represents an agent that accesses resources. It can represent an anonymous user, an
// authenticated user, or an internal Sourcegraph service.
type Actor struct {
	// UID is the unique ID of the authenticated user, or 0 for anonymous actors.
	UID int32 `json:",omitempty"`

	// Internal is true if the actor represents an internal Sourcegraph service (and is therefore
	// not tied to a specific user).
	Internal bool `json:",omitempty"`

	// FromSessionCookie is whether a session cookie was used to authenticate the actor. It is used
	// to selectivity display a logout link. (If the actor wasn't authenticated with a session 
	// cookie, logout would be ineffective.)
	FromSessionCookie bool `json:"-"`
}

func FromUser(uid int32) *Actor { return &Actor{UID: uid} }

type key int

const (
	actorKey key = iota
)

func WithActor(ctx context.Context, a *Actor) context.Context {
	if a != nil && a.UID != 0 {
		trace.TraceUser(ctx, a.UID)
	}
	return context.WithValue(ctx, actorKey, a)
}
