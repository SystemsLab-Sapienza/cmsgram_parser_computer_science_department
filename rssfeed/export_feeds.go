package rssfeed

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func exportFeeds(filepath string) error {
	conn := Pool.Get()
	defer conn.Close()

	// Get the number of RSS feeds in the database
	nfeeds, err := redis.Int(conn.Do("GET", "rss:feed:counter"))
	if err != nil {
		return err
	}

	// Create a new file for exporting
	file, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
		return err
	}

	// Create a new CSV writer for the file
	w := csv.NewWriter(file)

	for i := 1; i <= nfeeds; i++ {
		// Skip feeds set to be ignored
		ignore, err := redis.Bool(conn.Do("SISMEMBER", "rss:feed:ignore", i))
		if err != nil {
			return err
		}
		if ignore {
			continue
		}

		record, err := redis.Strings(conn.Do("HMGET", "rss:feed:"+strconv.Itoa(i), "name", "url", "kind"))
		if err != nil {
			log.Println(err)
			return err
		}

		if err = w.Write(record); err != nil {
			log.Println(err)
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
