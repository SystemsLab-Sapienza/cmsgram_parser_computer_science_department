package main

import (
	"bitbucket.org/ansijax/rfidlab_telegramdi_parser/config"
	// "bitbucket.org/ansijax/rfidlab_telegramdi_parser/newscrawler"
	"bitbucket.org/ansijax/rfidlab_telegramdi_parser/rssfeed"
)

func main() {
	var conf config.Config

	// Read the configuration from the file
	conf.Read()

	// Start the cralwer for the website
	// go newscrawler.Start()

	// Start the cralwer for the RSS feeds
	rssfeed.Start(conf)
}
