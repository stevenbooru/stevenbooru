package models

import "github.com/jinzhu/gorm"

// Tag is a tag that an image can have.
type Tag struct {
	gorm.Model
	UUID        string `sql:"unique;size:36" json:"uuid"`
	Name        string `sql:"unique;size:50" json:"name"`
	Description string `sql:"size:150" json:"description"`
}

type ImageTag struct {
	gorm.Model
	ImageID  uint
	TagID    uint
	SetterID uint
}
