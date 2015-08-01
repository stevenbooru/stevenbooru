package models

import (
	"net/http"

	"github.com/jinzhu/gorm"
)

// Image is an image that has been uploaded to the booru.
type Image struct {
	gorm.Model
	UUID        string `sql:"unique;size:36" json:"uuid"`
	Poster      *User
	PosterID    int
	Tags        []*Tag
	Hash        string `sql:"unique;size:128"`
	Filename    string `sql:"size:512"`
	Description string `sql:"size:2048"`
}

func NewImage(r *http.Request, user *User) (i *Image, err error) {
	return
}
