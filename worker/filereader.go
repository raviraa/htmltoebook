package worker

import (
	"io/ioutil"
	"os"
	"strings"
)

func ReadLines(fname string) []string {
	f, err := os.Open(fname)
	panicerr(err)
	b, err := ioutil.ReadAll(f)
	panicerr(err)
	links := SplitLinks(string(b))
	return links
}

func SplitLinks(s string) []string {
	// strip out invalid links, whitespaces
	splits := strings.Split(s, "\n")
	var links []string
	for _, s := range splits {
		s = strings.TrimSpace(s)
		if s != "" {

			links = append(links, s)
		}
	}
	return links
}
