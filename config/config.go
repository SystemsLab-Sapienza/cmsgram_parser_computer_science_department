package config

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	BotURI string

	RedisDomain      string
	RedisAddress     string
	RedisMaxIdle     int
	RedisIdleTimeout time.Duration
}

// Sets the default configuration
func (c *Config) Init() {
	c.BotURI = "/bot"

	c.RedisDomain = "tcp"
	c.RedisAddress = "localhost:6379"
	c.RedisMaxIdle = 3
	c.RedisIdleTimeout = 240
}

func (c *Config) Read(filepath string) error {
	c.Init()
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Can't find configuration file.", err)
	}

	r := csv.NewReader(f)
	r.Comma = ':'
	r.Comment = '#'
	r.FieldsPerRecord = 2
	r.LazyQuotes = true
	r.TrimLeadingSpace = true

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("main:", err)
		}

		value := record[1]
		switch record[0] {
		case "bot_URI":
			c.BotURI = value
		case "redis_domain":
			c.RedisDomain = value
		case "redis_address":
			c.RedisAddress = value
		case "redis_max_idle":
			i, err := strconv.Atoi(value)
			if err != nil {
				fmt.Printf("redis_max_idle value '%s' not valid. Using default.\n", value)
			} else {
				c.RedisMaxIdle = i
			}
		case "redis_idle_timeout":
			i, err := strconv.Atoi(value)
			if err != nil {
				fmt.Printf("redis_idle_timeout value '%s' not valid. Using default.\n", value)
			} else {
				c.RedisIdleTimeout = time.Duration(i)
			}
		default:
			fmt.Printf("Parameter '%s' in config file not valid. Ignored.\n", record[0])
		}
	}

	fmt.Printf("Server started with the following configuration:\n%-20s\t%s\n%-20s\t%s\n%-20s\t%s\n",
		"bot URI", c.BotURI,
		"redis domain:", c.RedisDomain,
		"redis address:", c.RedisAddress,
	)

	return err
}
