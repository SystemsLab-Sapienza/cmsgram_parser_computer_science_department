package rssfeed

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func ImportFeeds(filepath string) error {
	var feed = struct {
		Name string `redis:"name"`
		Kind string `redis:"kind"`
		Url  string `redis:"url"`
	}{}

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal("main: os.Open:", err) // TODO handle file not found
		return err
	}

	conn := Pool.Get()
	defer conn.Close()

	r := csv.NewReader(f)
	r.Comma = ','
	r.Comment = '#'
	r.FieldsPerRecord = 3
	r.LazyQuotes = false
	r.TrimLeadingSpace = true

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("main:", err)
		}

		feed.Name = record[0]
		feed.Url = record[1]
		feed.Kind = record[2]

		conn := Pool.Get()
		defer conn.Close()

		// Get current feed ID
		i, err := redis.Int(conn.Do("INCR", "rss:feed:counter"))
		if err != nil {
			return err
		}

		_, err = conn.Do("HMSET", redis.Args{}.Add("rss:feed:"+strconv.Itoa(i)).AddFlat(&feed)...)
		if err != nil {
			return err
		}
	}

	return nil
}
