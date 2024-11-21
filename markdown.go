package ktw

import (
	"context"
	"io"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type Markdown []byte

// Render Markdown into HTML.
//
// TODO: use github.com/abhinav/goldmark-frontmatter to parse the
// frontmatter and pass that back to the callee, so they can use it
// to run an appropriate template with the frontmatter variables as
// input.
func (m Markdown) Render(ctx context.Context, w io.Writer) error {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Strikethrough,
			extension.Table,
			extension.TaskList,
			extension.Typographer,
			// TODO: Install D2 extension... it seems github.com/FurqanSoftware/goldmark-d2
			// has an error with the included/used slog package. The authors must have some
			// sort of workspace override. As I'm not quite in need of the D2 extension, I'll
			// leave it as a todo at this point. Note, this would be as a means to surplant
			// my current use of Excalidraw+.
			// &d2.Extender{
			// 	Layout:  d2elklayout.Layout,
			// 	ThemeID: &d2themescatalog.Terminal.ID,
			// },
			NewCustomCodeHighlight(),
		),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
			parser.WithBlockParsers(
				util.Prioritized(NewInfoBlockParser(), 1000),
				util.Prioritized(NewNoteBlockParser(), 1000),
				util.Prioritized(NewWarningBlockParser(), 1000),
			),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithUnsafe(),
		),
	)
	return md.Convert(m, w)
}

// CustomCodeHighlight implements a custom fenced code highlighter
// that uses Chroma under the hood, but also implements class and
// general HTML attribute handling.
type CustomCodeHighlight struct {
	goldmark.Extender
}

// NewCustomCodeHighlight returns a wrapped Chroma highlight extension.
func NewCustomCodeHighlight() goldmark.Extender {
	wrapper := func(w util.BufWriter, context highlighting.CodeBlockContext, entering bool) {
		if entering {
			w.WriteString(`<pre class="chroma">`)
			lang, ok := context.Language()
			if !ok {
				lang = []byte("unknown")
			}
			// No attributes, write simple option
			if context.Attributes() == nil {
				w.WriteString(`<code class="language-`)
				w.Write(lang)
				w.WriteString(`">`)
				return
			}

			// Handle code with attributes
			w.WriteString(`<code class="language-`)
			w.Write(lang)
			if attr, ok := context.Attributes().GetString("class"); ok {
				w.WriteString(` `)
				w.Write(attr.([]byte))
				w.WriteString(`"`)
			}
			// Add in all the other possible attributes...
			for _, attr := range context.Attributes().All() {
				if !gmhtml.CodeAttributeFilter.Contains(attr.Name) {
					continue
				}
				if string(attr.Name) == "class" {
					continue
				}
				w.WriteString(` `)
				w.Write(attr.Name)
				w.WriteString(`="`)
				w.Write(attr.Value.([]byte))
				w.WriteString(`"`)
			}
			w.WriteString(`>`) // close out <code>
		} else {
			w.WriteString(`</code></pre>`)
		}
	}

	return &CustomCodeHighlight{
		Extender: highlighting.NewHighlighting(
			highlighting.WithFormatOptions(
				chromahtml.TabWidth(4),
				chromahtml.WithClasses(true),
				chromahtml.WithLineNumbers(true),
				chromahtml.WithPreWrapper(preWrapper{}),
			),
			highlighting.WithWrapperRenderer(wrapper),
		),
	}
}

type preWrapper struct{}

func (p preWrapper) Start(code bool, styleAttr string) string { return "" }

func (p preWrapper) End(code bool) string { return "" }

// Interface guard.
var _ Renderer = (*Markdown)(nil)
var _ chromahtml.PreWrapper = (*preWrapper)(nil)
