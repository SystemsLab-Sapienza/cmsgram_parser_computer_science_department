package rssfeed

import (
	"log"
	"os"
	"time"

	"infext/config"

	"github.com/garyburd/redigo/redis"
)

// Global variables
var (
	Config *config.Config
	Pool   *redis.Pool
)

func Start(conf *config.Config) {
	const delay = 1
	var flist FeedList

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

	return
}
