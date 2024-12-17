package ktw

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// InfoBlockParser implements a paragraph block parser that recognizes a
// paragraph starting with "Info:".
type InfoBlockParser struct {
	parser.BlockParser
}

// NewInfoBlockParser returns a BlockParser that recognizes a "Info:"
// start, which it then ensures that the class attribute of the resulting
// HTML element has a "note" class added to aid in styling.
func NewInfoBlockParser() parser.BlockParser {
	return &InfoBlockParser{
		BlockParser: parser.NewParagraphParser(),
	}
}

// Trigger will be triggered for lines starting with "Info:"
func (b *InfoBlockParser) Trigger() []byte {
	return []byte("Info:")
}

func (b *InfoBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	buf, _ := reader.PeekLine()
	found := strings.HasPrefix(string(buf), "Info:")
	pos := pc.BlockOffset()
	if pos < 0 || !found {
		return nil, parser.NoChildren
	}
	p, state := b.BlockParser.Open(parent, reader, pc)
	if p != nil {
		// Set the "class" attribute for the paragraph
		p.SetAttributeString("class", "info")
	}
	return p, state
}
