package rssfeed

import (
	"log"
	"os"
	"time"

	"bitbucket.org/ansijax/rfidlab_telegramdi_parser/config"

	"github.com/garyburd/redigo/redis"
)

// Global variables
var (
	Config config.Config
	Pool   *redis.Pool
)

func Start(conf config.Config) {
	var (
		delay = time.Duration(conf.RSSFeedDelay) // Delay in minutes
		flist FeedList
	)

	// If the goroutine panics at any point, don't bring down the whole program
	defer func() {
		if err := recover(); err != nil {
			log.Println("rssfeed panicked:", err)
		}
	}()

	// Create the log file
	f, err := os.OpenFile("rssfeed.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Can't create log file")
		return
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("Server started")

	// Create a thread-safe connection pool for redis
	Pool = &redis.Pool{
		MaxIdle:     conf.RedisMaxIdle,
		IdleTimeout: conf.RedisIdleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(conf.RedisDomain, conf.RedisAddress)
			if err != nil {
				return nil, err
			}

			return c, err
		},
	}

	Config = conf

	// Fetch the feed's data
	if err = loadFeeds(&flist); err != nil {
		log.Fatal("Error while loading feeds:", err)
	}

	// Start the main loop
	log.Println("Loading complete, entering main loop...")
	for {
		for _, f := range flist.feeds {
			feed := f
			go updateFeed(&feed)
		}
		time.Sleep(delay * time.Minute)
	}
}
