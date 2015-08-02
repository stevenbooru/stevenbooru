package main

import (
	"log"
	"time"

	"github.com/disintegration/imaging"
	"github.com/garyburd/redigo/redis"
	. "stevenbooru.cf/globals"
	"stevenbooru.cf/models"
)

func main() {
	handleStuff()
}

func handleStuff() {
	conn := Redis.Get()

	for {
		log.Println("Waiting on new image...")

		start := time.Now()

		ids, err := redis.Strings(conn.Do("BLPOP", "uploads", 0))
		if err != nil {
			doError(err)
			continue
		}

		id := ids[1]

		log.Printf("Found ID %v", ids)
		time.Sleep(2 * time.Second)

		if id == "" {
			log.Println("No new uploads found")
			time.Sleep(5 * time.Second)
		}

		img := &models.Image{}
		query := Db.Where("uuid = ?", id).First(img)
		if query.Error != nil {
			defer conn.Do("RPUSH", "uploads", id)
			doError(query.Error)
			continue
		}

		log.Printf("Starting to process image %d:%s", img.ID, img.UUID)

		image, err := imaging.Open(Config.Storage.Path + "/" + img.UUID + "/" + img.Filename)
		if err != nil {
			defer conn.Do("RPUSH", "uploads", id)
			doError(err)
			continue
		}

		croppedImage := imaging.CropCenter(image, 256, 256)
		if err != nil {
			defer conn.Do("RPUSH", "uploads", id)
			doError(err)
			continue
		}

		err = imaging.Save(croppedImage, Config.Storage.Path+"/"+img.UUID+"/thumbnail.jpg")
		if err != nil {
			defer conn.Do("RPUSH", "uploads", id)
			doError(err)
			continue
		}

		log.Printf("Wrote thumbnail for image %d:%s in %s", img.ID, img.UUID, time.Since(start).String())
	}

	conn.Close()
}

func doError(err error) {
	if Config.Site.Testing {
		log.Printf("%#v", err)

		time.Sleep(5 * time.Second)
	} else {
		panic(err)
	}
}
