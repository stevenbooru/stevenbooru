package models

import "github.com/jinzhu/gorm"

// User is a user on the Booru.
type User struct {
	gorm.Model
	UUID        string `sql:"size:36" json:"uuid"`         // UUID used in searches, etc
	ActualName  string `sql:"unique,size:75" json:"-"`     // lower case, unique name used in storage to prevent collisions
	DisplayName string `sql:"size:75" json:"display_name"` // user name that is displayed to users
	Email       string `sql:"size:400" json:"-"`           // email address for the user
	Role        string `json:"role"`                       // role that the user has on the booru

	PasswordHash string `json:"-"` // Blake2b hashed password of the user
	Salt         string `json:"-"` // Random data added to the password, along with the site's pepper

	// Relationships go here
}
