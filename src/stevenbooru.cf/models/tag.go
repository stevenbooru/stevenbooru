package models

import (
	"log"
	"strings"

	"github.com/Xe/uuid"
	"github.com/jinzhu/gorm"
	. "stevenbooru.cf/globals"
)

// Tag is a tag that an image can have.
type Tag struct {
	gorm.Model
	UUID        string `sql:"unique;size:36" json:"uuid"`
	Name        string `sql:"unique" json:"name"`
	Description string `sql:"size:150" json:"description"`
}

func NewTag(name string) (*Tag, error) {
	tag := &Tag{}
	var err error

	name = strings.ToLower(name)

	q := Db.Where(&Tag{
		Name: strings.ToLower(name),
	}).First(tag)
	if q.Error != nil {
		if q.Error.Error() != "record not found" {
			return nil, q.Error
		}
	}

	if tag.UUID == "" {
		log.Printf("Need to make a new tag for %s", name)

		tag = &Tag{
			UUID: uuid.New(),
			Name: name,
		}

		q := Db.Create(tag)
		if q.Error != nil {
			return nil, q.Error
		}
	}

	return tag, err
}

type ImageTag struct {
	gorm.Model
	ImageID  uint
	TagID    uint
	SetterID uint
}
