package parser

import (
	"fmt"
	"strings"

	"github.com/lucasepe/yamldo/text"
)

const spaceWidth = "  "

type Fragment struct {
	path    string
	isKey   bool
	depth   int
	content string
}

func (b *Fragment) Path() string {
	return b.path
}

func (b *Fragment) IsKey() bool {
	return b.isKey
}

func (b *Fragment) Level() int {
	return b.depth
}

func (b *Fragment) String() string {
	spaces := depthWidth(b.depth)
	res := text.Indent(b.content, spaces)
	if b.isKey {
		res = fmt.Sprintf("%s:", res)
	}

	return res
}

// depthWidth convert depth to spaces for indentation
func depthWidth(depth int) string {
	var rs strings.Builder
	for i := 0; i < depth; i++ {
		rs.WriteString(spaceWidth)
	}
	return rs.String()
}
