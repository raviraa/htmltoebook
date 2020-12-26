package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"unicode"

	readability "github.com/go-shiori/go-readability"
)

// TODO write out filename, title to a file
func fetchStripUrls(urls []string) {
	for i, url := range urls {
		article, err := readability.FromURL(url, 30*time.Second)
		// article, err := readability.FromReader(resp.Body, url)
		if err != nil {
			log.Fatalf("failed to parse %s, %v\n", url, err)
		}

		dstTxtFile, _ := os.Create(fmt.Sprintf("text-%02d.txt", i+1))
		defer dstTxtFile.Close()
		dstTxtFile.WriteString(article.TextContent)

		dstHTMLFile, _ := os.Create(fmt.Sprintf("html-%02d.html", i+1))
		defer dstHTMLFile.Close()
		dstHTMLFile.WriteString(article.Content)

		fmt.Printf("URL     : %s\n", url)
		fmt.Printf("Title   : %s\n", article.Title)
		fmt.Printf("Length  : %d\n", article.Length)
		fmt.Printf("Excerpt : %s\n", article.Excerpt)
		fmt.Printf("Text content saved to \"text-%02d.txt\"\n", i+1)
		fmt.Printf("HTML content saved to \"html-%02d.html\"\n", i+1)
		fmt.Println()
	}
}

// strips special charactors from a url
func urlToFname(url string) string {
	var out []rune
	for _, ch := range url {
		if unicode.IsLetter(ch) || unicode.IsNumber(ch) {
			out = append(out, ch)
		}
	}
	return string(out) + ".html"
}
