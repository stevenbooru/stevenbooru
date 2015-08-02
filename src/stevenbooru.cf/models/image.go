package models

import (
	"log"
	"net/http"
	"strings"

	"github.com/Xe/uuid"
	"github.com/jinzhu/gorm"
)

// Image is an image that has been uploaded to the booru.
type Image struct {
	gorm.Model
	UUID        string `sql:"unique;size:36" json:"uuid"`
	Poster      *User
	PosterID    uint
	Tags        []*Tag
	Hash        string `sql:"unique;size:128"`
	Filename    string `sql:"size:512"`
	Description string `sql:"size:2048"`
	Mime        string
}

func NewImage(r *http.Request, user *User) (i *Image, err error) {
	iuuid := uuid.NewUUID()

	_, header, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	tags := r.FormValue("tags")

	log.Printf("%#v", header.Header)
	log.Printf("%s", strings.Split(tags, ","))

	mime := header.Header.Get("Content-Type")

	i = &Image{
		UUID:     iuuid.String(),
		Poster:   user,
		PosterID: user.ID,
		Mime:     mime,
	}

	log.Printf("%#v", i)

	return
}
