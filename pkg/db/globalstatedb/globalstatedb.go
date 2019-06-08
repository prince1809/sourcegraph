package globalstatedb

import "context"

type State struct {
	SiteID      string
	Initialized bool // whether the initial site admin account has been created
}

func Get(ctx context.Context) (*State, error) {
	panic("")
}

func SiteInitialized(ctx context.Context) (alreadyInitialized bool, err error) {
	panic("")
}

func EnsureInitialized(ctx context.Context, dbh interface {
}) (alreadyInitialized bool, err error) {
	panic("")
}

func getConfiguration(ctx context.Context) (*State, error) {
	panic("")
}

func tryInsertNew(ctx context.Context, dbh interface {
}) error {
	panic("")
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
