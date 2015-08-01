package globals

import (
	"flag"
	"fmt"
	"log"
	"runtime"

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

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Set up the database
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

	// Turn on verbose logs for debugging
	if Config.Site.Testing {
		Db.LogMode(true)
	}

	// Test database connection
	err = Db.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database")

	// Set up redis
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

	// Test redis
	_, err = conn.Do("PING")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Redis")

	CookieStore = cookiestore.New([]byte(Config.Site.CookieHash))

	uuid.SetNodeID([]byte(Config.Site.Name))
}
