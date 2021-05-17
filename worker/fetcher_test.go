package worker

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/raviraa/htmltoebook/config"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestStripHtml(t *testing.T) {
	testhtml, err := os.Open("testdata/testhtml.html")
	require.Nil(t, err)
	defer testhtml.Close()
	w := Worker{conf: &config.ConfigType{AddPreBreaks: true}}
	b := new(bytes.Buffer)

	_, err = w.cleanHTML(testhtml, b, "http://localhost", nil)
	require.Nil(t, err)
	out := b.String()
	require.NotContains(t, out, "contains interface{}")
	require.Contains(t, out, "main() {<br/")
	require.Contains(t, out, "localhost/test.png")
	require.NotContains(t, out, `("script message")`)
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
	out := addPreDivBreaks(doc)
	require.Contains(t, out, "two(){<br/>")
	require.Contains(t, out, "main(){<br/>")
}

func TestDownloadImages(t *testing.T) {
	testhtml, err := os.Open("testdata/testhtml.html")
	require.Nil(t, err)
	defer testhtml.Close()
	w := New(nil, &config.ConfigType{Tmpdir: "."})

	b := new(bytes.Buffer)
	doc, err := html.Parse(testhtml)
	require.Nil(t, err)
	// TODO use httptest server
	out := w.downloadImages(doc, b)
	require.Contains(t, out, "../images/0001footer-gopher.jpg")
	require.Equal(t, "./0001footer-gopher.jpg\x000001footer-gopher.jpg\x00ADDIMAGE\n", b.String())
}
