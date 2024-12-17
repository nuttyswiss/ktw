package ktw

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// WarningBlockParser implements a paragraph block parser that recognizes a
// paragraph starting with "Warning:".
type WarningBlockParser struct {
	parser.BlockParser
}

// NewWarningBlockParser returns a BlockParser that recognizes a "Warning:"
// start, which it then ensures that the class attribute of the resulting
// HTML element has a "note" class added to aid in styling.
func NewWarningBlockParser() parser.BlockParser {
	return &WarningBlockParser{
		BlockParser: parser.NewParagraphParser(),
	}
}

// Trigger will be triggered for lines starting with "Warning:"
func (b *WarningBlockParser) Trigger() []byte {
	return []byte("Warning:")
}

func (b *WarningBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	buf, _ := reader.PeekLine()
	found := strings.HasPrefix(string(buf), "Warning:")
	pos := pc.BlockOffset()
	if pos < 0 || !found {
		return nil, parser.NoChildren
	}
	p, state := b.BlockParser.Open(parent, reader, pc)
	if p != nil {
		// Set the "class" attribute for the paragraph
		p.SetAttributeString("class", "warning")
	}
	return p, state
}
