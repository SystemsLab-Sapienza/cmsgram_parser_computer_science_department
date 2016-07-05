package main

import (
	"infext/config"
	// "infext/newscrawler"
	"infext/rssfeed"
)

func main() {
	var conf config.Config

	// Read the configuration from the file
	conf.Read()

	// Start the cralwer for the website
	// go newscrawler.Start()

	// Start the cralwer for the RSS feeds
	rssfeed.Start(&conf)
}
