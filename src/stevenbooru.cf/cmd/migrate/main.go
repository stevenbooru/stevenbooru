package main

import (
	"log"

	"stevenbooru.cf/globals"
	"stevenbooru.cf/models"
)

func main() {
	globals.Db.AutoMigrate(&models.User{})

	log.Println("Migration complete.")
}
