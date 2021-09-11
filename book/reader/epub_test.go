package reader

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEpub(t *testing.T) {
	b, err := NewEpub("testdata/book.epub")
	require.Nil(t, err)
	require.NotNil(t, b)

	out, err := os.Create("/tmp/book.txt")
	require.Nil(t, err)

	b.ReadFiles(func(sectionData []byte) {
		// txt, err := HtmltoTxt(sectionData)
		// require.Nil(t, err)
		txt := string(sectionData)
		fmt.Fprint(out, txt)
	})
	out.Close()
}
