package worker

import (
	"bytes"
	"os"
	"testing"

	"github.com/raviraa/htmltoebook/config"
	"github.com/stretchr/testify/require"
)

func TestStripHtml(t *testing.T) {
	testhtml, err := os.Open("testdata/testhtml.html")
	require.Nil(t, err)
	defer testhtml.Close()
	w := Worker{conf: &config.ConfigType{AddPreBreaks: true}}
	b := new(bytes.Buffer)

	_, err = w.stripHTML(testhtml, b, "http://localhost")
	require.Nil(t, err)
	out := b.String()
	require.NotContains(t, out, "contains interface{}")
	require.Contains(t, out, "main() {<br/")
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
	out := addPreDivBreaks(htm)
	require.Contains(t, out, "two(){<br/>")
	require.Contains(t, out, "main(){<br/>")
}
