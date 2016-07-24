package newscrawler

import (
	"crypto/md5"
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("unix", "/tmp/redis.sock")
			if err != nil {
				return nil, err
			}

			return c, err
		},
	}
}

/* This assumes a news is uniquely identified by its URL+publication date.
 */
func Start() {
	const delay = 5 // Delay in minutes

	// If the goroutine panics at any point, don't bring down the whole program
	defer func() {
		if err := recover(); err != nil {
			log.Println("newscrawler panicked:", err)
		}
	}()

	for {
		// Parse the website for news
		news, err := fetchNews()
		if err != nil {
			log.Fatal(err)
		}

		// Check result for new items
		// storeNewItems()

		conn := pool.Get()
		defer conn.Close()

		for _, n := range *news {
			checksum := md5.Sum([]byte(n.Title + n.Date))
			cs := string(checksum[:])

			exists, err := redis.Bool(conn.Do("SISMEMBER", "crawler:news", cs))
			if err != nil {
				log.Println(err)
				return
			}

			if exists {
				continue
			}

			ID, err := redis.Int(conn.Do("INCR", "crawler:news:counter"))
			if err != nil {
				log.Println(err)
				return
			}

			// Store the news on the DB
			conn.Send("MULTI")
			conn.Send("HMSET", redis.Args{}.Add("crawler:news:"+strconv.Itoa(ID)).AddFlat(&n)...)
			conn.Send("SADD", "crawler:news", cs)
			_, err = conn.Do("EXEC")
			if err != nil {
				log.Println(err)
				return
			}

			// Send update to bot
			go sendUpdate(strconv.Itoa(ID))
		}

		time.Sleep(delay * time.Minute)
	}

}
