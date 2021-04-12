package worker

import (
	"context"
	"fmt"
	"localhost/htmltoebook/config"
	"log"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/go-shiori/go-readability"
)

func FetchStripUrls(ctx context.Context, urls []string) bool {
	titlesfile, err := os.OpenFile(config.TitlesFname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer titlesfile.Close()

	for i, url := range urls {
		dstfname := config.Tmpdir + "/" + urlToFname(url)
		if _, err := os.Stat(dstfname); err == nil {
			loginfo("Ignoring cached url ", url)
			continue
		}
		loginfo(fmt.Sprintf("Fetching link %d/%d: %v", i+1, len(urls), url))

		resp, err := fetchUrl(ctx, url)
		if err != nil {
			logerr(url, err.Error())
			if config.Config.FailonError {
				log.Fatal(err)
			}
			continue
		}
		defer resp.Body.Close()

		article, err := readability.FromReader(resp.Body, url)
		if err != nil {
			logerr("failed to parse ", url, err.Error())
			continue
		}

		dstHTMLFile, err := os.Create(dstfname)
		if err != nil {
			log.Fatal(err)
		}
		defer dstHTMLFile.Close()
		dstHTMLFile.WriteString(article.Content)

		fmt.Fprintf(titlesfile, "%s %s\n", dstfname, article.Title)

		loginfo(fmt.Sprintf("Sleeping for %d seconds ", config.Config.SleepSec))
		// time.Sleep(time.Duration(config.Config.SleepSec * int(time.Second)))
		select {
		case <-time.After(time.Duration(config.Config.SleepSec * int(time.Second))):
			// pass on normal timeout
		case <-ctx.Done():
			logerr("Stopping the process")
			return false
		}
	}
	return true
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

func fetchUrl(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", config.Config.UserAgent)
	req = req.WithContext(ctx)

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
