package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/go-shiori/go-readability"
)

func (w *Worker) FetchStripUrls(ctx context.Context, urls []string) bool {
	titlesfile, err := os.OpenFile(w.conf.TitlesFname(), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		w.logerr("unable to open intermediate files for writing ", err.Error())
		return false
	}
	defer titlesfile.Close()

	for i, url := range urls {
		dstfname := w.conf.Tmpdir + "/" + urlToFname(url)
		if _, err := os.Stat(dstfname); err == nil {
			w.loginfo("Ignoring cached url ", url)
			continue
		}
		w.loginfo(fmt.Sprintf("Fetching link %d/%d: %v", i+1, len(urls), url))

		resp, err := w.fetchUrl(ctx, url)
		if err != nil {
			w.logerr(url, err.Error())
			if w.conf.FailonError {
				// log.Fatal(err)
				w.logerr(err.Error())
				return false
			}
			continue
		}
		defer resp.Body.Close()

		article, err := readability.FromReader(resp.Body, url)
		if err != nil {
			w.logerr("failed to parse ", url, err.Error())
			continue
		}

		dstHTMLFile, err := os.Create(dstfname)
		if err != nil {
			log.Fatal(err)
		}
		defer dstHTMLFile.Close()
		dstHTMLFile.WriteString(article.Content)

		fmt.Fprintf(titlesfile, "%s\x00%s\n", dstfname, article.Title)

		w.loginfo(fmt.Sprintf("Fetched article with %d letters. Sleeping for %d seconds ", article.Length, w.conf.SleepSec))
		select {
		case <-time.After(time.Duration(w.conf.SleepSec * int(time.Second))):
			// nothing to do on normal timeout
		case <-ctx.Done():
			w.logerr("Stopping the process")
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

func (w *Worker) fetchUrl(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", w.conf.UserAgent)
	req = req.WithContext(ctx)

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("Error fetching url " + resp.Status)
		return nil, err
	}
	return resp, nil
}
