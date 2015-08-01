package models

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/Xe/uuid"
	"github.com/dchest/blake2b"
	"github.com/jinzhu/gorm"
	. "stevenbooru.cf/globals"
)

var (
	ErrUserCreateMissingValues = errors.New("models.User: missing values on creation attempt")
	ErrInvalidEmail            = errors.New("models: bad email address")
	ErrDifferentPasswords      = errors.New("models: the same password was not used twice")
)

// User is a user on the Booru.
type User struct {
	gorm.Model
	UUID        string `sql:"unique;size:36" json:"uuid"`  // UUID used in searches, etc
	ActualName  string `sql:"unique;size:75" json:"-"`     // lower case, unique name used in storage to prevent collisions
	DisplayName string `sql:"size:75" json:"display_name"` // user name that is displayed to users
	Email       string `sql:"size:400" json:"-"`           // email address for the user
	Role        string `json:"role"`                       // role that the user has on the booru
	AvatarURL   string `json:"avatar_url"`                 // URL to the user's avatar
	Activated   bool   `json:"-"`                          // Has the user activated their email address?

	PasswordHash string `json:"-"` // Blake2b hashed password of the user
	Salt         string `json:"-"` // Random data added to the password, along with the site's pepper

	// Relationships go here
}

// NewUser makes a new user in the database given the values from a HTTP POST request.
func NewUser(values url.Values) (u *User, err error) {
	username := values.Get("username")
	if username == "" {
		return nil, ErrUserCreateMissingValues
	}

	email := values.Get("email")
	if email == "" {
		return nil, ErrUserCreateMissingValues
	}

	// TODO: check for duplicate email addresses
	if !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}

	password := values.Get("password")
	if password == "" {
		return nil, ErrUserCreateMissingValues
	}

	confirm := values.Get("password_confirm")
	if confirm == "" {
		return nil, ErrUserCreateMissingValues
	}

	if password != confirm {
		return nil, ErrDifferentPasswords
	}

	salt := uuid.New()[0:14]

	c := &blake2b.Config{
		Salt: []byte(salt),
		Size: blake2b.Size,
	}

	b2b, err := blake2b.New(c)
	if err != nil {
		panic(err)
	}

	b2b.Reset()

	fin := bytes.NewBufferString(password + Config.Site.Pepper)
	io.Copy(b2b, fin)

	result := fmt.Sprintf("%x", b2b.Sum(nil))

	myUuid := uuid.NewUUID().String()

	u = &User{
		Email:        email,
		DisplayName:  username,
		ActualName:   url.QueryEscape(strings.ToLower(username)),
		Activated:    false,
		UUID:         myUuid,
		Salt:         salt,
		PasswordHash: result,
	}

	Db.Create(u)

	if Db.NewRecord(u) {
		return nil, errors.New("could not create user")
	}

	return
}
