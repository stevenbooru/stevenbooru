package models

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Xe/uuid"
	"github.com/dchest/blake2b"
	"github.com/jinzhu/gorm"
	. "stevenbooru.cf/globals"
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
	userUpload, header, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(userUpload)
	if err != nil {
		return nil, err
	}

	c := &blake2b.Config{
		Salt: []byte("lol as if"),
		Size: blake2b.Size,
	}

	b2b, err := blake2b.New(c)
	if err != nil {
		panic(err)
	}

	b2b.Reset()
	io.Copy(b2b, bytes.NewBuffer(data))
	hash := fmt.Sprintf("%x", b2b.Sum(nil))

	tags := r.FormValue("tags")

	log.Printf("%#v", header.Header)
	log.Printf("%s", strings.Split(tags, ","))

	mime := header.Header.Get("Content-Type")

	i = &Image{
		UUID:        uuid.New(),
		Poster:      user,
		PosterID:    user.ID,
		Mime:        mime,
		Filename:    header.Filename,
		Hash:        hash,
		Description: r.Form.Get("description"),
	}

	err = os.Mkdir(Config.Storage.Path+"/"+i.UUID, os.ModePerm)
	if err != nil {
		return nil, err
	}

	fout, err := os.Create(Config.Storage.Path + "/" + i.UUID + "/" + i.Filename)
	if err != nil {
		return nil, err
	}

	io.Copy(fout, bytes.NewBuffer(data))

	conn := Redis.Get()
	defer conn.Close()

	_, err = conn.Do("RPUSH", "uploads", i.UUID)
	if err != nil {
		return nil, err
	}

	query := Db.Create(i)
	if query.Error != nil {
		return nil, query.Error
	}

	return
}
