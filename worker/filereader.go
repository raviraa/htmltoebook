package worker

import (
	"io/ioutil"
	"os"
	"strings"
)

func ReadLines(fname string) ([]string, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	links := SplitLinks(string(b))
	return links, nil
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
