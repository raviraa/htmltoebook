package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	nurl "net/url"
	"os"
	"path"
	"strings"
	"time"
	"unicode"

	"github.com/go-shiori/dom"
	"github.com/go-shiori/go-readability"
	"github.com/raviraa/htmltoebook/config"
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
		dstfname := w.conf.Tmpdir() + "/" + urlToFname(url)
		if _, err := os.Stat(dstfname); err == nil {
			w.loginfo("Ignoring cached url ", url)
			if w.conf.ChromeDownload && w.browser == nil {
				// chrome is started only each ChromeRestartNum. this will start it in case of cached and not mod first
				w.resetChrome()
			}
			continue
		}
		w.loginfo(fmt.Sprintf("Fetching link %d/%d: %v", i+1, len(urls), url))

		if i%config.ChromeRestartNum == 0 {
			w.resetChrome()
		}
		var timeStart time.Time
		var htmbody string
		// retry failed requests RetryNum times, restarting chrome if needed. TODO listen for ctx.cancel
		for i := 0; i < config.RetryNum; i++ {
			timeStart = time.Now()
			htmbody, err = w.fetchBody(ctx, url)
			if err == nil {
				goto FetchEnd
			} else {
				w.logerr(fmt.Sprint("try: ", i, ", ", url, err.Error()))
				if i == config.RetryNum-1 {
					if w.conf.FailonError {
						// log.Fatal(err)
						w.logerr(err.Error())
						return false
					}
					goto FetchCont
				}

				w.resetChrome()
				log.Println("Sleeping for ", w.conf.SleepSec)
				time.Sleep(time.Duration(w.conf.SleepSec * int(time.Second)))
			}
		}
	FetchCont: // ignore fetch error and continue with next link
		continue
	FetchEnd: // fetch success.

		dstHTMLFile, err := os.Create(dstfname)
		if err != nil {
			log.Fatal(err)
		}
		article, err := w.cleanHTML(htmbody, dstHTMLFile, url, titlesfile)
		if err != nil {
			continue
		}
		dstHTMLFile.Close()

		if w.conf.FailonError && article.Length < w.conf.MinPageSize {
			os.WriteFile(
				path.Join(os.TempDir(), "emptyhtm"),
				[]byte(htmbody), 0755)
			os.WriteFile(
				path.Join(os.TempDir(), "emptycontent"),
				[]byte(article.Content), 0755)
			w.logerr(fmt.Sprintf("Article length: %v less than required MinPageSize: %v ", article.Length, w.conf.MinPageSize))
			os.Remove(dstfname)
			return false
		}

		fmt.Fprintf(titlesfile, "%s\x00%s\x00%s\n", dstfname, urlToFname(article.Title), url)

		sleeprand := w.conf.SleepSec/2 + rand.Intn(w.conf.SleepSec/2)
		w.loginfo(fmt.Sprintf("Fetched article with %d characters in %v. Sleeping for %d seconds ", article.Length, time.Since(timeStart), sleeprand))
		select {
		case <-time.After(time.Duration(sleeprand * int(time.Second))):
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
		if unicode.IsLetter(ch) || unicode.IsNumber(ch) || ch == ' ' {
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
	// req.Header.Add("Accept-Encoding", "gzip, deflate")
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
func (w *Worker) cleanHTML(body string, out io.Writer, url string, titlesfile io.Writer) (readability.Article, error) {
	r := bytes.NewBufferString(body)
	netUrl, _ := nurl.ParseRequestURI(url)
	article, err := readability.FromReader(r, netUrl)
	if err != nil {
		w.logerr("failed to parse ", url, err.Error())
		return article, err
	}

	// Useful for epub pages with parsing issues"`
	if w.conf.UseTextParser {
		fmt.Fprintf(out, "<h4>%s</h4>", article.Title)
		htm := article.TextContent
		htm = strings.ReplaceAll(htm, "\n", "<br>")
		out.Write([]byte(htm))
		return article, nil
	}

	htm := article.Content
	htmb := []byte(htm)
	if w.conf.PreBreaks || w.conf.IncludeImages {
		doc, err := html.Parse(strings.NewReader(htm))
		if err != nil {
			log.Fatal(err)
		} else {
			if w.conf.PreBreaks {
				htmb = addPreDivBreaks(doc)
			}
			if w.conf.IncludeImages {
				htmb = w.downloadImages(doc, titlesfile)
			}
		}
	}

	// remove extra html>body tags added by html parser
	htmb = bytes.TrimPrefix(htmb, []byte("<html><head></head><body>"))
	htmb = bytes.TrimSuffix(htmb, []byte("</body></html>"))
	// add title
	out.Write([]byte(fmt.Sprintf("<h3>%s</h3>\n", article.Title)))
	out.Write(htmb)
	return article, nil
}

// adds <br> at the end of each line of <pre> block
func addPreDivBreaks(doc *html.Node) []byte {

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
	return b.Bytes()
}

// downloads images in doc(DOM), and updates src attribute to relative file path(../images/img.png)
func (w *Worker) downloadImages(doc *html.Node, titlesfile io.Writer) []byte {
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
				if ctype != "image/png" && ctype != "image/jpg" && ctype != "image/jpeg" {
					continue
				}
				if len(ctype) > 6 && ctype[:6] == "image/" {
					imgext = "." + ctype[6:]
				}
				imgfname := fmt.Sprintf("%04d%s%s", w.imgCount, path.Base(attr.Val), imgext)
				dstimgfname := w.conf.Tmpdir() + "/" + imgfname
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
	return b.Bytes()
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

func (w *Worker) fetchBody(ctx context.Context, url string) (string, error) {
	if w.conf.ChromeDownload {
		return w.FetchUrlChrome(url)
	}

	resp, err := w.fetchUrl(ctx, url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	htmbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(htmbody), nil
}
