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
	"path"
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
		article, err := w.cleanHTML(resp.Body, dstHTMLFile, url, titlesfile)
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
	if titlesfile.Close() != nil {
		w.logerr("error writing file. " + err.Error())
		return false
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

// converts to readable html and fetches images if necessary
func (w *Worker) cleanHTML(r io.ReadCloser, out io.Writer, url string, titlesfile io.Writer) (readability.Article, error) {

	parser := readability.NewParser()
	article, err := parser.Parse(r, url)
	if err != nil {
		w.logerr("failed to parse ", url, err.Error())
		return article, err
	}
	r.Close()

	htm := article.Content
	doc, err := html.Parse(strings.NewReader(htm))
	if err != nil {
		log.Println(err)
	} else {
		if w.conf.PreBreaks {
			htm = addPreDivBreaks(doc)
		}
		if w.conf.IncludeImages {
			htm = w.downloadImages(doc, titlesfile)
		}

	}

	out.Write([]byte(htm))
	return article, nil
}

// adds <br> at the end of each line of <pre> block
func addPreDivBreaks(doc *html.Node) string {

	preNodes := dom.GetElementsByTagName(doc, "pre")
	for _, node := range preNodes {
		nodeTxt := dom.TextContent(node)
		nodeTxt = strings.ReplaceAll(nodeTxt, "\n", "<br>\n")
		dom.SetInnerHTML(node, nodeTxt)
	}

	b := new(bytes.Buffer)
	err := html.Render(b, doc)
	if err != nil {
		log.Println(err)
	}
	return b.String()
}

// downloads images in doc(DOM), and updates src attribute to relative file path(../images/img.png)
func (w *Worker) downloadImages(doc *html.Node, titlesfile io.Writer) string {
	preNodes := dom.GetElementsByTagName(doc, "img")
	for _, node := range preNodes {
		for _, attr := range node.Attr {
			if attr.Key == "src" {
				w.loginfo("fetching image " + attr.Val)
				// TODO get correct context for cancel
				resp, err := w.fetchUrl(context.Background(), attr.Val)
				if err != nil {
					w.logerr("error fetching image " + err.Error())
					continue
				}
				w.imgCount++ // counter for unique image file name
				// set file extension from Content-Type
				imgext := ""
				ctype := resp.Header.Get("Content-Type")
				if len(ctype) > 6 && ctype[:6] == "image/" {
					imgext = "." + ctype[6:]
				}
				imgfname := fmt.Sprintf("%04d%s%s", w.imgCount, path.Base(attr.Val), imgext)
				dstimgfname := w.conf.Tmpdir + "/" + imgfname
				if err = writefile(dstimgfname, resp.Body); err != nil {
					w.logerr("error fetching image " + err.Error())
					continue
				}
				dom.SetAttribute(node, "src", fmt.Sprintf("../images/%s", imgfname))
				fmt.Fprintf(titlesfile, "%s\x00%s\x00%s\n", dstimgfname, imgfname, ADDIMAGE)
				// time.Sleep(time.Duration(w.conf.SleepSec) * time.Second)
			}
		}
	}
	b := new(bytes.Buffer)
	err := html.Render(b, doc)
	if err != nil {
		log.Println(err)
	}
	return b.String()
}

func writefile(fname string, r io.ReadCloser) error {
	defer r.Close()
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}
