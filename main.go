package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	flag.StringVar(&config.LinksFile, "l", "", "file containing http links to fetch (REQUIRED)")
	flag.BoolVar(&config.FailonError, "f", false, "exit when fetching any of the link fails")
	flag.StringVar(&config.UserAgent, "ua", "Mozilla/5.0 (X11; rv:84.0) Gecko/20100101 Firefox/84.0", "user agent to use for http request")
	flag.IntVar(&config.SleepInterval, "s", 3, "sleep interval between fetching links")

	flag.Parse()
	if config.LinksFile == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	links := readLines(config.LinksFile)
	fetchStripUrls(links)
	writeMobi()

}

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}

// log to websocket in future?
func out(s ...string) {
	log.Println(s)
}
