package main

import (
	"fmt"
	"log"
	"os"

	"github.com/raviraa/htmltoebook/cmd"
	"github.com/raviraa/htmltoebook/web"
)

const usage = `Usage: htmltoebook [mode]
Optional mode can be one of the following:
 (w|web)		Starts web interface in browser (default)
 (c|console)		Starts console mode urls editor with $EDITOR. Can also edit settings.
 (s|snippet)		Starts console mode snippets editor with $EDITOR. Parses urls from given html code snippet.
`

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	var mode string
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "w", "web", "":
		web.NewWeb()
	case "c", "console":
		cmd.RunLinks()
	case "s", "snippet":
		cmd.RunHtmlSnippet()
	default:
		fmt.Println(usage)
		os.Exit(2)
	}
}
