package worker

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/raviraa/htmltoebook/config"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestStripHtml(t *testing.T) {
	testhtml, err := ioutil.ReadFile("testdata/testhtml.html")
	require.Nil(t, err)
	w := Worker{conf: &config.ConfigType{PreBreaks: true}}
	b := new(bytes.Buffer)

	_, err = w.cleanHTML(string(testhtml), b, "http://localhost", nil)
	require.Nil(t, err)
	out := b.String()
	require.NotContains(t, out, "contains interface{}")
	require.Contains(t, out, "main() {<br/")
	require.Contains(t, out, "localhost/test.png")
	require.NotContains(t, out, `("script message")`)
	require.NotContains(t, out, "<html>")
	require.Contains(t, out, "<h3>Test File Title</h3")
}

func TestAddPreBreaks(t *testing.T) {
	htm := `
	<div>test one</div>
	<pre>
	import "fmt"
	func main(){

	}
	</pre>
	<div>test two</div>
	<pre>
	import "fmt"
	func two(){
	}
	</pre>

	`
	doc, err := html.Parse(strings.NewReader(htm))
	require.Nil(t, err)
	out := string(addPreDivBreaks(doc))
	require.Contains(t, out, "two(){<br/>")
	require.Contains(t, out, "main(){<br/>")
}

func TestDownloadImages(t *testing.T) {
	testhtml, _ := os.Open("testdata/testhtml.html")
	testhtm, _ := ioutil.ReadAll(testhtml)
	testhtml.Close()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "image/png")
		f, _ := os.Open("../screenshot.png")
		io.Copy(w, f)
	}))
	w := New(nil, nil, &config.ConfigType{DownloadDir: "."})
	htm := strings.ReplaceAll(string(testhtm), "SRCIMG", ts.URL+"/test")
	b := new(bytes.Buffer)

	doc, err := html.Parse(strings.NewReader(htm))
	require.Nil(t, err)
	out := string(w.downloadImages(doc, b))
	require.Contains(t, out, "../images/0001test.png")
	require.Equal(t, "/tmp/0001test.png\x000001test.png\x00ADDIMAGE\n", b.String())
	os.Remove("0001test.png")
}

func TestStripHtmlPara(t *testing.T) {
	testhtml, err := ioutil.ReadFile("testdata/testhtmlpara.html")
	require.Nil(t, err)
	w := Worker{conf: &config.ConfigType{UseTextParser: true}}
	b := new(bytes.Buffer)

	_, err = w.cleanHTML(string(testhtml), b, "http://localhost", nil)
	require.Nil(t, err)
	require.Contains(t, b.String(), "one<br>")
}
