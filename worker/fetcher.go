package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/go-shiori/dom"
	"github.com/go-shiori/go-readability"
	"golang.org/x/net/html"
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

		dstHTMLFile, err := os.Create(dstfname)
		if err != nil {
			log.Fatal(err)
		}
		defer dstHTMLFile.Close()
		article, err := w.stripHTML(resp.Body, dstHTMLFile, url)
		if err != nil {
			continue
		}

		fmt.Fprintf(titlesfile, "%s\x00%s\x00%s\n", dstfname, article.Title, url)

		w.loginfo(fmt.Sprintf("Fetched article with %d characters. Sleeping for %d seconds ", article.Length, w.conf.SleepSec))
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

// strips special characters from a url
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

func (w *Worker) stripHTML(r io.Reader, out io.Writer, url string) (readability.Article, error) {

	parser := readability.NewParser()
	article, err := parser.Parse(r, url)
	if err != nil {
		w.logerr("failed to parse ", url, err.Error())
		return article, err
	}

	htm := article.Content
	if w.conf.AddPreBreaks {
		htm = addPreDivBreaks(htm)
	}

	out.Write([]byte(htm))
	return article, nil
}

// adds <br> at the end of each line of <pre> block
func addPreDivBreaks(str string) string {
	doc, err := html.Parse(strings.NewReader(str))
	if err != nil {
		log.Println(err)
		return str
	}

	preNodes := dom.GetElementsByTagName(doc, "pre")
	for _, node := range preNodes {
		nodeTxt := dom.TextContent(node)
		nodeTxt = strings.ReplaceAll(nodeTxt, "\n", "<br>\n")
		dom.SetInnerHTML(node, nodeTxt)
	}

	b := new(bytes.Buffer)
	err = html.Render(b, doc)
	if err != nil {
		log.Println(err)
	}
	return b.String()
}
