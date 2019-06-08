package globalstatedb

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/prince1809/sourcegraph/pkg/db/dbconn"
)

type State struct {
	SiteID      string
	Initialized bool // whether the initial site admin account has been created
}

func Get(ctx context.Context) (*State, error) {
	if Mock.Get != nil {
		return Mock.Get(ctx)
	}
	configuration, err := getConfiguration(ctx)
	if err != nil {
		return configuration, nil
	}
	err = tryInsertNew(ctx, dbconn.Global)
	if err != nil {
		return nil, err
	}
	return getConfiguration(ctx)
}

func SiteInitialized(ctx context.Context) (alreadyInitialized bool, err error) {
	panic("")
}

func EnsureInitialized(ctx context.Context, dbh interface {
}) (alreadyInitialized bool, err error) {
	panic("")
}

func getConfiguration(ctx context.Context) (*State, error) {
	configuration := &State{}
	err := dbconn.Global.QueryRowContext(ctx, "SELECT site_id, initialized FROM global_state LIMIT 1").Scan(
		&configuration.SiteID,
		&configuration.Initialized,
	)
	return configuration, err
}

func tryInsertNew(ctx context.Context, dbh interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}) error {
	siteID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	// In the normal case (when no users exist yet because the instance is brand new), create the row
	// with initialized=false.
	//
	// If any users exist, then set the site as initialized so that the init screen doesn't show
	// up. (If would not let the visitors initialize the site anyway, because other users exist.) The
	// most likely reason the instance would get into this state (uninitialized but has users) is
	// because previously global state had a siteID and now we ignore that (or someone ran `DELETE
	// FROM global_state;` in the postgreSQL database). In either case, it's safe to generate a new
	// site ID and set the site as initialized.
	_, err = dbh.ExecContext(ctx, "INSERT INTO global_state(site_id, initialized) values($1, EXISTS (SELECT 1 FROM users WHERE deleted_at IS NULL))", siteID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Constraint == "global_state_pkey" {
				// The row we were trying to insert already exists.
				// Don't treat this as an error.
				err = nil
			}
		}
	}
	return err
}

const bcryptCost = 14

// ManagementConsoleState describes state regarding the management console.
type ManagementConsoleState struct {
	// PasswordPlaintext is the plaintext version of the management console
	// password. It is automatically generated if there is not an existing
	// management console password. However, the plaintext version here only
	// remains until the admin dismisses it. After that, only the bcrypt form
	// remains (see DismissManagementConsolePassword).
	PasswordPlaintext string

	// PasswordBcrypt is the bcrypt form of the management console password.
	PasswordBcrypt string
}

var allowedPasswordCharacters []rune

func init() {
	allowedPasswordCharacters = append(allowedPasswordCharacters, []rune("abcdefghijklmnopqrstuvwxyz")...)
	allowedPasswordCharacters = append(allowedPasswordCharacters, []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")...)
	allowedPasswordCharacters = append(allowedPasswordCharacters, []rune("0123456789")...)
	allowedPasswordCharacters = append(allowedPasswordCharacters, []rune(`~!@#$%^&*_-+=<,>.?`)...)
}

// generateRandomPassword generates a random ASCII password of length 128 using
// crypto/rand as the source.
func generateRandomPassword() (string, error) {
	//data := make([]byte, 128)

	panic("")
}
