package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/raviraa/htmltoebook/cmd"
	"github.com/raviraa/htmltoebook/web"
)

func usage() {
	fmt.Print(`Usage: htmltoebook [flags] [mode]
Optional mode can be one of the following:
 (w|web)		Starts web interface in browser
 (c|console)		Starts console mode urls editor with $EDITOR. Can also edit settings. (default)
 (s|snippet)		Starts console mode snippets editor with $EDITOR. Parses urls from given html code snippet.

`)
	flag.PrintDefaults()
}

func main() {
	var mode string
	log.SetFlags(log.Lshortfile | log.Ltime)
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) == 1 {
		mode = flag.Args()[0]
	}

	switch mode {
	case "w", "web":
		web.NewWeb()
	case "c", "console", "": // default mode
		cmd.RunLinks("")
	case "s", "snippet":
		cmd.RunHtmlSnippet()
	default:
		usage()
		os.Exit(2)
	}
}
