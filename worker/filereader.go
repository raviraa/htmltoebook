package worker

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
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
	links := strings.Split(string(b), "\n")
	return links, nil
}

// expands http://two.com/page{{1-3}}/content to 3 links
var rangeLinks = regexp.MustCompile(`http.*{{(.*)}}.*`)
var rangeLinksReplace = regexp.MustCompile(`({{.*}})`)

func SplitLinks(s string) []string {
	// strip out invalid links, whitespaces
	splits := strings.Split(s, "\n")
	var links []string
	for _, s := range splits {
		s = strings.TrimSpace(s)
		if len(s) > 4 && s[:4] == "http" {
			if match := rangeLinks.FindStringSubmatch(s); len(match) > 1 {
				var numlo, numhi int
				if n, _ := fmt.Sscanf(match[1], "%d-%d", &numlo, &numhi); n == 2 {
					for i := numlo; i <= numhi; i++ {
						link := rangeLinksReplace.ReplaceAllString(s, strconv.Itoa(i))
						links = append(links, link)
					}
					continue
				}
			}
			links = append(links, s)
		}
	}
	return links
}
