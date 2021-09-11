package writer

import "github.com/766b/mobi"

type mobiWriter struct {
	book *mobi.MobiWriter
}

func NewMobi(title, fname string) BookWriter {
	m, err := mobi.NewWriter(fname)
	if err != nil {
		return nil
	}

	m.Title(title)
	m.Compression(mobi.CompressionNone) // LZ77 compression is also possible using  mobi.CompressionPalmDoc
	// m.Compression(mobi.CompressionPalmDoc)
	// Meta data
	// m.NewExthRecord(mobi.EXTH_DOCTYPE, "EBOK")
	// m.NewExthRecord(mobi.EXTH_AUTHOR, "Book Author Name")

	return &mobiWriter{m}
}

func (e *mobiWriter) AddSection(title, body string) {
	// e.book.AddSection(body, title, "", "")
	// m.NewChapter(titles[fname], b)
	e.book.NewChapter(title, []byte(body))
}

func (e *mobiWriter) Write() error {
	e.book.Write()
	return nil
}

func (e *mobiWriter) AddImage(imagefile, imageattrsrc string) error {
	return nil
}

/*


	}
	// Output MOBI File
	m.Write()


*/
