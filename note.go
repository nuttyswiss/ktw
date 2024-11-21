package ktw

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// NoteBlockParser implements a paragraph block parser that recognizes a
// paragraph starting with "Note:".
type NoteBlockParser struct {
	parser.BlockParser
}

// NewNoteBlockParser returns a BlockParser that recognizes a "Note:"
// start, which it then ensures that the class attribute of the resulting
// HTML element has a "note" class added to aid in styling.
func NewNoteBlockParser() parser.BlockParser {
	return &NoteBlockParser{
		BlockParser: parser.NewParagraphParser(),
	}
}

// Trigger will be triggered for lines starting with "Note:"
func (b *NoteBlockParser) Trigger() []byte {
	return []byte("Note:")
}

func (b *NoteBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	p, state := b.BlockParser.Open(parent, reader, pc)
	if p != nil {
		// Set the "class" attribute for the paragraph
		p.SetAttributeString("class", "note")
	}
	return p, state
}
