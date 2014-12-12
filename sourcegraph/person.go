package sourcegraph

import (
	"crypto/md5"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"

	"sourcegraph.com/sourcegraph/go-sourcegraph/db_common"

	"sourcegraph.com/sourcegraph/go-nnz/nnz"
)

// User represents a user.
type User struct {
	// UID is the numeric primary key for a user.
	UID UID `db:"uid"`

	// GitHubID is the numeric ID of the GitHub user account corresponding to
	// this user.
	GitHubID nnz.Int `db:"github_id"`

	// Login is the user's username, which typically corresponds to the user's
	// GitHub login.
	Login string

	// Name is the (possibly empty) full name of the user.
	Name string

	// Type is either "User" or "Organization".
	Type string

	// AvatarURL is the URL to an avatar image specified by the user.
	AvatarURL string

	// Location is the user's physical location (from their GitHub profile).
	Location string `json:",omitempty"`

	// Company is the user's company (from their GitHub profile).
	Company string `json:",omitempty"`

	// HomepageURL is the user's homepage or blog URL (from their GitHub
	// profile).
	HomepageURL string `db:"homepage_url" json:",omitempty"`

	// Transient is if this user was constructed on the fly and is not persisted
	// or resolved to a Sourcegraph/GitHub/etc. user.
	Transient bool `db:"-" json:",omitempty"`

	// UserProfileDisabled is whether the user profile should not be displayed
	// on the Web app.
	UserProfileDisabled bool `db:"user_profile_disabled" json:",omitempty"`

	// RegisteredAt is the date that the user registered. If the user has not
	// registered (i.e., we have processed their repos but they haven't signed
	// into Sourcegraph), it is null.
	RegisteredAt db_common.NullTime `db:"registered_at"`
}

// GitHubLogin returns the user's Login. They are the same for now, but callers
// that intend to get the GitHub login should call GitHubLogin() so that we can
// decouple the logins in the future if needed.
func (u *User) GitHubLogin() string {
	return u.Login
}

// IsOrganization is whether this user represents a GitHub organization
// (which are treated as a subclass of User in GitHub's data model).
func (u *User) IsOrganization() bool { return u.Type == "Organization" }

func (u *User) AvatarURLOfSize(width int) string {
	return u.AvatarURL + fmt.Sprintf("&s=%d", width)
}

// CanOwnRepositories is whether the user is capable of owning repositories
// (e.g., GitHub users can own GitHub repositories).
func (u *User) CanOwnRepositories() bool {
	return u.GitHubLogin() != ""
}

// CanAttributeCodeTo is whether this user can commit code. It is false for
// organizations and true for both users and transient users.
func (u *User) CanAttributeCodeTo() bool {
	return !u.IsOrganization()
}

// CanSync is whether this person can be synced with the external source that
// the person was originally fetched from (e.g., GitHub users).
func (u *User) CanSync() bool {
	return !u.Transient
}

// UID is the numeric primary key for a user.
type UID int

// Scan implements database/sql.Scanner.
func (x *UID) Scan(v interface{}) error {
	if data, ok := v.(int64); ok {
		*x = UID(data)
		return nil
	}
	return fmt.Errorf("%T.Scan failed: %v", x, v)
}

// Value implements database/sql/driver.Valuer.
func (x UID) Value() (driver.Value, error) {
	return int64(x), nil
}

// DefaultAvatarSize is the size, in pixels, of avatar images if no size is
// specified.
const DefaultAvatarSize = 128

// GravatarURL returns the URL to the Gravatar avatar image for email. If size
// is 0, DefaultAvatarSize is used.
func GravatarURL(email string, size uint16) string {
	if size == 0 {
		size = DefaultAvatarSize
	}
	h := md5.New()
	io.WriteString(h, email)
	return fmt.Sprintf("https://secure.gravatar.com/avatar/%x?s=%d&d=mm", h.Sum(nil), size)
}

// ErrPersonNotExist is an error indicating that no such user exists.
var ErrPersonNotExist = errors.New("user does not exist")

// ErrPersonRenamed is an error type that indicates that a user account was renamed
// from OldLogin to NewLogin.
type ErrPersonRenamed struct {
	// OldLogin is the previous login name.
	OldLogin string

	// NewLogin is what the old login was renamed to.
	NewLogin string
}

func (e ErrPersonRenamed) Error() string {
	return fmt.Sprintf("login %q was renamed to %q; use the new name", e.OldLogin, e.NewLogin)
}

type PersonStatType string

type PersonStats map[PersonStatType]int

const (
	PersonStatAuthors            = "authors"
	PersonStatClients            = "clients"
	PersonStatOwnedRepos         = "owned-repos"
	PersonStatContributedToRepos = "contributed-to-repos"
	PersonStatDependencies       = "dependencies"
	PersonStatDependents         = "dependents"
	PersonStatDefs               = "defs"
	PersonStatExportedDefs       = "exported-defs"
)

func (x PersonStatType) Value() (driver.Value, error) {
	return string(x), nil
}

func (x *PersonStatType) Scan(v interface{}) error {
	if data, ok := v.([]byte); ok {
		*x = PersonStatType(data)
		return nil
	}
	return fmt.Errorf("%T.Scan failed: %v", x, v)
}
