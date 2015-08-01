package globals

import (
	"flag"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"stevenbooru.cf/config"
)

var (
	Config config.Config
	Db     gorm.DB

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
}
