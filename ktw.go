package ktw

import (
	"bytes"
	"context"
	"io"

	elem "github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
)

// Renderer implements the Render functionality of any components.
type Renderer interface {
	Render(context.Context, io.Writer) error
}

// Page represents a single Web page. It contains everything necessary to let
// the page's render function produce an HTML representation of the page and
// all its contents.
type Page struct {
	Title    string
	Metadata map[string]any
	Contents []Renderer
}

// Render produces the HTML representing this page and all its contents.
//
// TODO: Use frontmatter style to create a template to evaluate with the
// frontmatter as input to the template. Find the template files by
// evaluating the provided, and parsed, configuration file.
func (p *Page) Render(ctx context.Context, w io.Writer) error {
	var body []elem.Node
	for _, item := range p.Contents {
		var raw bytes.Buffer
		if err := item.Render(ctx, &raw); err != nil {
			return err
		}
		body = append(body, elem.Raw(raw.String()))
	}

	html := elem.Html(nil,
		elem.Head(nil,
			elem.Meta(attrs.Props{attrs.Charset: "utf-8"}),
			elem.Title(nil, elem.Text(p.Title)),
			elem.Comment("Generated by Magic"),
		),
		elem.Body(nil, body...),
	)
	_, err := w.Write([]byte(html.Render()))
	return err
}

// Interface guard.
var _ Renderer = (*Page)(nil)
