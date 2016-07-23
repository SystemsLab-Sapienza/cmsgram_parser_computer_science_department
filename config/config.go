package config

import (
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
	*c = Config{"http://localhost:8443/wl08ncvrqisnv1wu8unwl08k05vo81j9", "unix", "/tmp/redis.sock", 3, 240}
}

func (c *Config) Read() {
	c.Init()
}
