package models

import "github.com/jinzhu/gorm"

// Image is an image that has been uploaded to the booru.
type Image struct {
	gorm.Model
	UUID     string `sql:"unique;size:36" json:"uuid"`
	Poster   *User
	PosterID int
	Tags     []*Tag
	Hash     string `sql:"unique;size:128"`
}
