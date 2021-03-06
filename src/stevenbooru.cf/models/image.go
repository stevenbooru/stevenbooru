package models

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Xe/uuid"
	"github.com/dchest/blake2b"
	"github.com/jinzhu/gorm"
	. "stevenbooru.cf/globals"
)

var (
	ErrNeedRatingTag = errors.New("models: need a rating tag to proceed")
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

	tagsRaw := r.FormValue("tags")
	tags := strings.Split(tagsRaw, ",")

	for _, tag := range tags {
		if strings.HasPrefix(tag, "rating:") {
			goto ok
		}
	}

	return nil, ErrNeedRatingTag

ok:

	mime := header.Header.Get("Content-Type")

	if mime == "image/svg+xml" {
		return nil, errors.New("Unsupported image format")
	}

	if !strings.HasPrefix(mime, "image/") {
		return nil, errors.New("Unsupported image format")
	}

	i = &Image{
		UUID:        uuid.New(),
		Poster:      user,
		PosterID:    user.ID,
		Mime:        mime,
		Filename:    header.Filename,
		Hash:        hash,
		Description: r.Form.Get("description"),
	}

	for _, tag := range tags {
		t, err := NewTag(tag)
		if err != nil {
			return nil, err
		}

		i.Tags = append(i.Tags, t)
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

	query := Db.Create(i)
	if query.Error != nil {
		return nil, query.Error
	}

	_, err = conn.Do("RPUSH", "uploads", i.UUID)
	if err != nil {
		return nil, err
	}

	return
}
