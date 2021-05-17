package htmlextract

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseHtmlLinks(t *testing.T) {
	htm := `
        <form id="form2" live-submit="evsnippetsave">
        <div>
          <label for="HtmlSnippet">Html code snippet to extract anchor links from</label>
					<a class="foo" href="/pageone">first link</a>
					</div>
					<p><a href="http://li2.com">two link wont match</a>
					<p><a href="http://li3.com">third link</a>

					<a>dummy without link</a>
	`
	links, err := ParseHtmlLinks(htm, "http://one.com", "", "", false)
	require.Nil(t, err)
	require.Equal(t, 3, len(links))
	require.Equal(t, "http://one.com/pageone", links[0])

	links, err = ParseHtmlLinks(htm, "", "", `li\d.com`, false)
	require.Nil(t, err)
	require.Equal(t, 2, len(links))

	links, err = ParseHtmlLinks(htm, "", "link$", "", false)
	require.Nil(t, err)
	require.Equal(t, 2, len(links))
	htm = `
        <form id="form2" live-submit="evsnippetsave">
        <div>
          <label for="HtmlSnippet">Html code snippet to extract anchor links from</label>
	`
	links, _ = ParseHtmlLinks(htm, "", "", "", false)
	require.Equal(t, 0, len(links))
}
