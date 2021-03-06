package models

import (
	"errors"
	"net/url"
	"strings"

	"github.com/Xe/uuid"
	"github.com/jinzhu/gorm"
	. "stevenbooru.cf/globals"
)

var (
	ErrUserCreateMissingValues = errors.New("models.User: missing values on creation attempt")
	ErrInvalidEmail            = errors.New("models: bad email address")
	ErrDifferentPasswords      = errors.New("models: the same password was not used twice")
	ErrBadPassword             = errors.New("models: user gave an incorrect password")
)

// User is a user on the Booru.
type User struct {
	gorm.Model
	UUID        string `sql:"unique;size:36" json:"uuid"`  // UUID used in searches, etc
	ActualName  string `sql:"unique;size:75" json:"-"`     // lower case, unique name used in storage to prevent collisions
	DisplayName string `sql:"size:75" json:"display_name"` // user name that is displayed to users
	Email       string `sql:"unique;size:400" json:"-"`    // email address for the user
	Role        string `json:"role"`                       // role that the user has on the booru
	AvatarURL   string `json:"avatar_url"`                 // URL to the user's avatar
	Activated   bool   `json:"-"`                          // Has the user activated their email address?

	PasswordHash string `json:"-"` // Blake2b hashed password of the user
	Salt         string `json:"-"` // Random data added to the password, along with the site's pepper
}

func Login(values url.Values) (u *User, err error) {
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

	user := &User{}

	query := Db.Where("email = ?", email).First(&user)
	if query.Error != nil {
		return nil, query.Error
	}

	if user.checkPassword(password) == false {
		return nil, ErrBadPassword
	}

	return user, nil
}

func (u *User) checkPassword(password string) bool {
	return u.PasswordHash == hashPassword(password, u.Salt)
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
	result := hashPassword(password, salt)

	myUuid := uuid.NewUUID().String()

	u = &User{
		Email:        email,
		DisplayName:  username,
		ActualName:   url.QueryEscape(strings.ToLower(username)),
		Activated:    false,
		UUID:         myUuid,
		Salt:         salt,
		PasswordHash: result,
		AvatarURL:    "/img/avatar_default.jpg",
	}

	query := Db.Create(u)
	if query.Error != nil {
		return nil, query.Error
	}

	if Db.NewRecord(u) {
		return nil, errors.New("could not create user")
	}

	return
}
