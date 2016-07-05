package rssfeed

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func sendUpdate(f *Feed) error {
	var payload = struct {
		Key   string
		Value string
	}{"rss", strconv.Itoa(f.ID)}

	// Encode the payload into JSON
	data, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	// Send the payload TODO get URL from config file
	_, err = http.Post(Config.BotUri, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	return nil
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
		_, err = conn.Do("HSET", feed, "last_update", fdate) // TODO
		if err != nil {
			log.Println(err)
			return
		}

		if err = sendUpdate(f); err != nil {
			log.Println(err)
			return
		}
	}

	return
}
