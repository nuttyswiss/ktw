package ktw

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func md(doc string) Markdown {
	return Markdown(strings.ReplaceAll(doc, "'", "`"))
}

func TestHtmlRenderer(t *testing.T) {
	ctx := context.Background()
	pg := &Page{
		Title:    "Test Title",
		Contents: []Renderer{md(testdoc1)},
	}
	var buf bytes.Buffer

	err := pg.Render(ctx, &buf)
	if err != nil {
		t.Errorf("Got error: %+v", err)
	}
	t.Logf("Got:\n%s", buf.String())
}

func TestMarkdownRender(t *testing.T) {
	ctx := context.Background()
	md := md(testdoc1)
	var buf bytes.Buffer

	err := md.Render(ctx, &buf)
	if err != nil {
		t.Errorf("Got error: %+v", err)
	}
	t.Logf("Got:\n%s", buf.String())
}

var testdoc1 = `# header one{.note}

Note: This is a note!!
With extra lines, but still part of the same note.

This paragraph is not a note, but contains a
Note: at the start of a line.

Info: a quick informational note.

Warning: This is bad news! But we can have a bit of
inline code like 'rm -rf / && echo whoops'.

This is a quick paragraph with 'test -e /tmp && echo yes':

'''go {.bad .very-bad #example1}
package main

import "fmt"

func main() {
	fmt.Println("Hello World!")
}
'''{.good}
`
