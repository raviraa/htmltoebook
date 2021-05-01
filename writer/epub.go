package writer

import (
	"github.com/bmaupin/go-epub"
)

type epubWriter struct {
	book  *epub.Epub
	fname string
}

func NewEpub(title, fname string) BookWriter {
	return &epubWriter{
		book:  epub.NewEpub(title),
		fname: fname,
	}
}

func (e *epubWriter) AddSection(title, body string) {
	e.book.AddSection(body, title, "", "")
}

func (e *epubWriter) Write() error {
	return e.book.Write(e.fname)
}
