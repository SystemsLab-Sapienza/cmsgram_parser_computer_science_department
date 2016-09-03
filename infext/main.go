package main

import (
	"flag"
	"log"

	"bitbucket.org/ansijax/rfidlab_telegramdi_parser/config"
	"bitbucket.org/ansijax/rfidlab_telegramdi_parser/newscrawler"
	"bitbucket.org/ansijax/rfidlab_telegramdi_parser/rssfeed"
)

var flagConfigFile string

func init() {
	flag.StringVar(&flagConfigFile, "c", "", "Specifies the path to the config file.")
}

func main() {
	var conf config.Config

	flag.Parse()

	if len(flagConfigFile) == 0 {
		log.Fatal("You have to provide a configuration file: infext -c /path/to/file")
	}

	// Read the configuration from the file
	conf.Read(flagConfigFile)

	// Start the crawler for the website
	go newscrawler.Start(conf)

	// Start the crawler for the RSS feeds
	rssfeed.Start(conf)
}
