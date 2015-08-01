package globals

import (
	"flag"
	"fmt"
	"log"

	"github.com/Xe/uuid"
	"github.com/garyburd/redigo/redis"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"stevenbooru.cf/config"
)

var (
	Config      config.Config
	Db          gorm.DB
	Redis       *redis.Pool
	CookieStore sessions.Store

	ConfigFileFlag = flag.String("conf", "./cfg/stevenbooru.cfg", "configuration file to load")
	IrcConfigFlag  = flag.String("ircconf", "./cfg/irc.cfg", "config file for the IRC bots")
)

func init() {
	flag.Parse()

	var err error
	Config, err = config.ParseConfig(*ConfigFileFlag)
	if err != nil {
		log.Fatal(err)
	}

	Db, err = gorm.Open(Config.Database.Kind,
		fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
			Config.Database.Username,
			Config.Database.Password,
			Config.Database.Database,
			Config.Database.Host,
			Config.Database.Port,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = Db.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database")

	Redis = &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", Config.Redis.Host, Config.Redis.Port))
			if err != nil {
				return nil, err
			}

			if Config.Redis.Password != "" {
				if _, err := c.Do("AUTH", Config.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
	}

	conn := Redis.Get()
	defer conn.Close()

	_, err = conn.Do("PING")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Redis")

	CookieStore = cookiestore.New([]byte(Config.Site.CookieHash))

	uuid.SetNodeID([]byte(Config.Site.Name))
}
