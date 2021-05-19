package htmlextract

import (
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/html"
)

type ParseReq struct {
	HtmlSnippet   string `toml:"-"`
	BaseUrl       string `comment:"Url to use as base for relative links"`
	UrlNameRegex  string `comment:"Regex to match text name of the anchor tag\n Use single quotes if it contains escape sequences. ex.: '\\d+' "`
	UrlRegex      string `comment:"Regex to match href url part of anchor tag"`
	ExcludeFilter bool   `comment:"Matching anchor links will be excluded"`
	ReverseList   bool   `comment:"rever the list of links matched"`
}

// extract anchor links from raw html snippet
// filtered by regex for href and anchor text
func ParseHtmlLinks(r ParseReq) (links []string, err error) {
	log.Println("Parsing doc with length: ", len(r.HtmlSnippet))
	doc, err := html.Parse(strings.NewReader(r.HtmlSnippet))
	if err != nil {
		return
	}
	anchorRe, err := regexp.Compile(r.UrlNameRegex)
	if err != nil {
		return nil, err
	}
	linkRe, err := regexp.Compile(r.UrlRegex)
	if err != nil {
		return nil, err
	}
	aNodes := dom.GetElementsByTagName(doc, "a")
	for _, node := range aNodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				// fmt.Println(attr, dom.InnerText(node))
				match := linkRe.MatchString(attr.Val) && anchorRe.MatchString(dom.InnerText(node))
				if r.ExcludeFilter && match || !r.ExcludeFilter && !match {
					continue
				}
				if len(attr.Val) > 4 && attr.Val[:4] != "http" {
					links = append(links, r.BaseUrl+attr.Val)
				} else {
					links = append(links, attr.Val)
				}
			}
		}
	}
	if r.ReverseList {
		links = reverseSlice(links)
	}
	return links, nil
}

func NewParseReq(input map[string]interface{}) (*ParseReq, error) {
	var out ParseReq
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &out,
		// map from form sends boolean true as "on"
		DecodeHook: func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			// fmt.Println(f, f.Kind(), t.Kind())
			if f.Kind() == reflect.String && t.Kind() == reflect.Bool {
				return data.(string) == "on", nil
			}
			return data, nil
		},
	}
	dec, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}
	if err = dec.Decode(input); err != nil {
		return nil, err
	}
	return &out, nil
}

func reverseSlice(a []string) []string {
	for i, j := 0, len(a)-1; i < j; {
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}
	return a
}
