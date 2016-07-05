package rssfeed

import (
	"log"
	"strconv"
	"sync"

	"github.com/SlyMarbo/rss"
	"github.com/garyburd/redigo/redis"
)

func loadFeeds(flist *FeedList) error {
	var wg sync.WaitGroup

	conn := Pool.Get()
	defer conn.Close()

	// Fetch the number of feeds
	nfeeds, err := redis.Int(conn.Do("GET", "rss:feed:counter"))
	if err != nil {
		return err
	}

	for i := 1; i <= nfeeds; i++ {
		// Get the feed's URL
		furl, err := redis.String(conn.Do("HGET", "rss:feed:"+strconv.Itoa(i), "url"))
		if err != nil {
			return err
		}

		// Copy the value of i in a newly allocated variable accessible from the closure called in this iteration
		ID := i

		// Add a new goroutine to the wait group
		wg.Add(1)

		// Dispatch a closure in a new goroutine to fetch the feed and save the data
		// this way the requests are made concurrently and we skip right to the main loop
		go func() {
			defer wg.Done()
			log.Println("Fetching feed:", furl)
			feed, err := rss.Fetch(furl)
			if err != nil {
				log.Println("Can't fetch feed:", furl, err)
				return
			}

			// Atomically update the slice
			flist.Add(ID, feed)
			log.Println("Fetching complete:", furl)
		}()
	}

	// Wait for all goroutines to finish updating the feed list
	wg.Wait()

	return nil
}
