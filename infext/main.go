package main

import (
	"fmt"
	"os"

	"bitbucket.org/ansijax/rfidlab_telegramdi_parser/config"
	"bitbucket.org/ansijax/rfidlab_telegramdi_parser/newscrawler"
	"bitbucket.org/ansijax/rfidlab_telegramdi_parser/rssfeed"
)

func main() {
	var conf config.Config

	if len(os.Args) == 1 {
		fmt.Println("No config file provided, exiting.")
		return
	}

	// Read the configuration from the file
	conf.Read(os.Args[1])

	// Start the crawler for the website
	go newscrawler.Start(conf)

	// Start the crawler for the RSS feeds
	rssfeed.Start(conf)
}
