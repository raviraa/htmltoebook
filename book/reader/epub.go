package reader

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type epubBook struct {
	book *zip.ReadCloser
}

func NewEpub(fname string) (BookReader, error) {
	r, err := zip.OpenReader(fname)
	if err != nil {
		return nil, fmt.Errorf("unable to read ebook %s, %w", fname, err)
	}

	book := epubBook{r}
	return book, nil
}

func (e epubBook) ReadFiles(sectionData func([]byte)) error {
	for _, f := range e.book.File {
		fmt.Printf("Contents of %s:\n", f.Name)
		if !strings.HasPrefix(f.Name, "EPUB/xhtml/section") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			log.Println(err)
			return err
		}
		b, err := ioutil.ReadAll(rc)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("lennn", len(b))
		sectionData(b)
		rc.Close()
	}

	e.book.Close()
	return nil
}

/*
func HtmltoTxt(htm []byte) (out string, err error) {
	parser := readability.NewParser()
	article, err := parser.Parse(bytes.NewReader(htm), "/")
	if err != nil {
		return
	}
	out = fmt.Sprintf("\n\n### %s \n%s", article.Title, article.TextContent)
	return
}
*/
