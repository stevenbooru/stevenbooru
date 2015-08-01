package main

import (
	"log"

	"stevenbooru.cf/globals"
	"stevenbooru.cf/models"
)

func main() {
	globals.Db.AutoMigrate(&models.User{}, &models.Tag{}, &models.Image{})

	log.Println("Migration complete.")
}
