package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/766b/mobi"
)

func writeMobi() {
	m, err := mobi.NewWriter("output.mobi")
	if err != nil {
		panic(err)
	}

	m.Title("Book Title")
	// m.Compression(mobi.CompressionNone) // LZ77 compression is also possible using  mobi.CompressionPalmDoc
	m.Compression(mobi.CompressionPalmDoc)

	// Add cover image
	// m.AddCover("data/cover.jpg", "data/thumbnail.jpg")

	// Meta data
	m.NewExthRecord(mobi.EXTH_DOCTYPE, "EBOK")
	m.NewExthRecord(mobi.EXTH_AUTHOR, "Book Author Name")
	// See exth.go for additional EXTH record IDs

	// TODO title
	txtfiles, err := filepath.Glob("*html")
	chkerr(err)
	for _, txtfile := range txtfiles {
		fmt.Println("Adding ", txtfile)
		b, err := ioutil.ReadFile(txtfile)
		chkerr(err)

		// b, heading := txtToHtml(b)
		m.NewChapter(txtfile, b)
	}

	// Output MOBI File
	m.Write()

}

// func txtToHtml(b []byte) ([]byte, string) {
// 	out := bytes.NewBuffer(nil)
// 	heading := ""
// 	for _, bs := range bytes.Split(b, []byte("\n")) {
// 		if heading == "" {
// 			heading = string(bs)
// 		}
// 		s := fmt.Sprintf("<p>%s</p>", string(bs))
// 		out.WriteString(s)
// 	}

// 	return out.Bytes(), heading
// }
