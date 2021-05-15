package htmlextract

import (
	"regexp"
	"strings"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

// extract anchor links from raw html snippet
// filtered by regex for href and anchor text
func ParseHtmlLinks(htm, baseUrl, anchorRegex, linkRegex string) (links []string, err error) {
	doc, err := html.Parse(strings.NewReader(htm))
	if err != nil {
		return
	}
	anchorRe, err := regexp.Compile(anchorRegex)
	if err != nil {
		return nil, err
	}
	linkRe, err := regexp.Compile(linkRegex)
	if err != nil {
		return nil, err
	}
	aNodes := dom.GetElementsByTagName(doc, "a")
	for _, node := range aNodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				// fmt.Println(attr, dom.InnerText(node))
				if linkRe.MatchString(attr.Val) && anchorRe.MatchString(dom.InnerText(node)) {
					if len(attr.Val) > 0 && attr.Val[:4] != "http" {
						links = append(links, baseUrl+attr.Val)
					} else {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
	return links, nil
}
