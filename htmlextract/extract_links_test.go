package htmlextract

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	// links, err := ParseHtmlLinks(ParseReq{htm, "http://one.com", "", "", false})
	links, err := ParseHtmlLinks(ParseReq{HtmlSnippet: htm, BaseUrl: "http://one.com"})
	require.Nil(t, err)
	require.Equal(t, 3, len(links))
	require.Equal(t, "http://one.com/pageone", links[0])

	// links, err = ParseHtmlLinks(ParseReq{htm, "", "", `li\d.com`, false})
	links, err = ParseHtmlLinks(ParseReq{HtmlSnippet: htm, UrlRegex: `li\d.com`})
	require.Nil(t, err)
	require.Equal(t, 2, len(links))

	// links, err = ParseHtmlLinks(ParseReq{htm, "", "link$", "", false})
	links, err = ParseHtmlLinks(ParseReq{HtmlSnippet: htm, UrlNameRegex: `link$`})
	require.Nil(t, err)
	require.Equal(t, 2, len(links))
	htm = `
        <form id="form2" live-submit="evsnippetsave">
        <div>
          <label for="HtmlSnippet">Html code snippet to extract anchor links from</label>
	`
	links, err = ParseHtmlLinks(ParseReq{HtmlSnippet: htm})
	require.Equal(t, 0, len(links))
	require.Nil(t, err)
}

func TestNewParseReq(t *testing.T) {
	m := map[string]interface{}{
		"UrlRegex":      "UrlRegex",
		"BaseUrl":       "http://localhost",
		"ExcludeFilter": "on",
	}
	pr, err := NewParseReq(m)
	assert.Nil(t, err)
	assert.Equal(t, true, pr.ExcludeFilter)
	assert.Equal(t, "UrlRegex", pr.UrlRegex)
	assert.Equal(t, pr.BaseUrl, "http://localhost")
}
