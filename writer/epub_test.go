package writer

import (
	"os"
	"testing"

	"github.com/bmaupin/go-epub"
	"github.com/stretchr/testify/require"
)

func TestEpubWrite(t *testing.T) {
	// Create a new EPUB
	e := epub.NewEpub("My title")

	// Set the author
	e.SetAuthor("Hingle McCringleberry")

	// Add a section
	section1Body := `<html><body><h1>Section 1</h1>
<p>This is a paragraph.</p></body></html>`
	e.AddSection(section1Body, "Section 1", "", "")

	err := e.Write("/tmp/tt.epub")
	require.Nil(t, err)
	s, err := os.Stat("/tmp/tt.epub")
	require.Nil(t, err)
	require.Greater(t, s.Size(), int64(1000))
	os.Remove("/tmp/tt.epub")
}
