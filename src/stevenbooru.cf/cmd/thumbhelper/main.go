package main

import (
	"log"
	"runtime"
	"time"

	"github.com/disintegration/imaging"
	"github.com/garyburd/redigo/redis"
	. "stevenbooru.cf/globals"
	"stevenbooru.cf/models"
)

func main() {
	for i := 0; i < (runtime.NumCPU() - 1); i++ {
		go handleStuff()
	}

	handleStuff()
}

func handleStuff() {
	conn := Redis.Get()

	defer func() {
		if r := recover(); r != nil {
			handleStuff()
		}
	}()

	for {
		log.Println("Waiting on new image...")

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

		croppedImage := imaging.Thumbnail(image, 256, 256, imaging.Box)
		if err != nil {
			defer conn.Do("RPUSH", "uploads", id)
			doError(err)
			continue
		}

		err = imaging.Save(croppedImage, Config.Storage.Path+"/"+img.UUID+"/thumbnail.png")
		if err != nil {
			defer conn.Do("RPUSH", "uploads", id)
			doError(err)
			continue
		}

		log.Printf(
			"Wrote thumbnail for image %d:%s",
			img.ID,
			img.UUID,
		)
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
