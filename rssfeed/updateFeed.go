package rssfeed

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func sendUpdate(f *Feed) (err error) {
	var payload = struct {
		Key   string
		Value string
	}{"rss", strconv.Itoa(f.ID)}

	// Encode the payload into JSON
	data, err := json.Marshal(&payload)
	if err != nil {
		return
	}

	// Send the payload
	_, err = http.Post(Config.BotURI, "application/json", bytes.NewReader(data))
	if err != nil {
		return
	}

	return
}

func updateFeed(f *Feed) (err error) {
	feed := "rss:feed:" + strconv.Itoa(f.ID)

	conn := Pool.Get()
	defer conn.Close()

	err = f.Update()
	if err != nil {
		log.Println(err)
		return
	}

	lupdate, err := redis.Int64(conn.Do("HGET", feed, "last_update"))
	if err != nil && err != redis.ErrNil {
		log.Println(err)
		return
	}

	if len(f.Items) == 0 {
		return
	}

	fdate := f.Items[0].Date.Unix()
	if lupdate == 0 || fdate > lupdate {
		// Set the new time of last update
		_, err = conn.Do("HSET", feed, "last_update", fdate)
		if err != nil {
			log.Println(err)
			return
		}

		if lupdate != 0 {
			if err = sendUpdate(f); err != nil {
				log.Println(err)
				return
			}
		}
	}

	return
}
