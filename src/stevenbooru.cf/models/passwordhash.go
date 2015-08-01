package models

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dchest/blake2b"
	. "stevenbooru.cf/globals"
)

func hashPassword(password, salt string) (result string) {
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

	result = fmt.Sprintf("%x", b2b.Sum(nil))

	return
}

func CheckPassword(password, salt, hashed string) bool {
	return hashed == hashPassword(password, salt)
}
