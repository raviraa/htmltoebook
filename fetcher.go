package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/go-shiori/go-readability"
)

func fetchStripUrls(urls []string) {
	titlesfile, err := os.OpenFile(titlesfname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer titlesfile.Close()

	for i, url := range urls {
		dstfname := tmpdir + "/" + urlToFname(url)
		if _, err := os.Stat(dstfname); err == nil {
			out("Ignoring fetched url ", url)
			continue
		}

		out(fmt.Sprintf("Fetching link %d/%d", i+1, len(urls)))
		resp, err := fetchUrl(url)

		if err != nil {
			out(url, err.Error())
			if config.FailonError {
				panic(err)
			}
			continue
		}
		defer resp.Body.Close()

		article, err := readability.FromReader(resp.Body, url)
		if err != nil {
			out("failed to parse ", url, err.Error())
			continue
		}

		dstHTMLFile, err := os.Create(dstfname)
		if err != nil {
			panic(err)
		}
		defer dstHTMLFile.Close()
		dstHTMLFile.WriteString(article.Content)

		fmt.Fprintf(titlesfile, "%s %s\n", dstfname, article.Title)

		out(fmt.Sprintf("Sleeping for duration %d seconds ", config.SleepInterval))
		time.Sleep(time.Duration(config.SleepInterval * int(time.Second)))
	}
}

// strips special charactors from a url
func urlToFname(url string) string {
	var out []rune
	for _, ch := range url {
		if unicode.IsLetter(ch) || unicode.IsNumber(ch) {
			out = append(out, ch)
		}
	}
	return string(out) + ".html"
}

func fetchUrl(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", config.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
