package worker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitLinks(t *testing.T) {
	s := `
ignored line
http://two.com/page{{1-3}}/content
  http://one.com  `

	links := SplitLinks(s)
	assert.Equal(t, 4, len(links))
	assert.Equal(t, "http://two.com/page1/content", links[0])
	assert.Equal(t, "http://one.com", links[len(links)-1])
}
