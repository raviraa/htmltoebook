package writer

import (
	"log"
	"os"
	"testing"

	"github.com/bmaupin/go-epub"
	"github.com/stretchr/testify/require"
)

func TestEpubWrite(t *testing.T) {
	tmpfname := "/tmp/out.epub"
	e := epub.NewEpub("My title")

	// Set the author
	e.SetAuthor("Author name")

	// Add a section
	section1Body := `<html><body><h1>Section 1</h1>
<p>This is a paragraph.</p>
<p>And second paragraph with image <img src="../images/go-gopher.png" /> </p>
</body></html>`
	e.AddSection(section1Body, "Section 1", "", "")

	img1Path, err := e.AddImage("testdata/screenshot.png", "go-gopher.png")
	if err != nil {
		log.Fatal(err)
	}
	require.Nil(t, err)
	require.Equal(t, "../images/go-gopher.png", img1Path)

	err = e.Write(tmpfname)
	require.Nil(t, err)
	s, err := os.Stat(tmpfname)
	require.Nil(t, err)
	require.Greater(t, s.Size(), int64(1000))
	os.Remove(tmpfname)
}
