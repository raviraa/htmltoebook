package main

import (
	"flag"
	"localhost/htmltoebook/config"
	"localhost/htmltoebook/web"
	"localhost/htmltoebook/worker"
	"log"
	"os"
)

func main2() {

	flag.StringVar(&config.Config.LinksFile, "l", "", "file containing http links to fetch (REQUIRED)")
	flag.BoolVar(&config.Config.FailonError, "f", false, "exit when fetching any of the link fails")

	flag.Parse()
	if config.Config.LinksFile == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	links := worker.ReadLines(config.Config.LinksFile)
	worker.FetchStripUrls(links)
	worker.WriteMobi()

}

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	web.NewWeb()
	// webWorker(h)
}
