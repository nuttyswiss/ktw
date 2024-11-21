/*
Package ktw implements fuctionality to write Web applications in a Go centric
manner. It implements a simple Markdown converter with extra functionality.

It uses a customized Goldmark Markdown processor that contains a few enhancements.
In particular, it allows for special styling of "Note:", "Info:", and "Warning:"
paragraphs, which are paragraphs that start with the aforementioned words (with
the following ":"). These will be rendered as HTML "<p>" elements with a class of
"note", "info", and "warning" respectively.

The code fence has also been enhanced. It uses the Chroma extension to parse and
provide spans with classes (as before), but it adds the ability to add attributes
to the "<code>" element. In particular, it revives the "class=language-..." class,
and adds any other classes and/or IDs that you mention in a custom attribute
extension:

```go {.good #example1}
package awesome
```

Would result in a code element such as: <code class="language-go good" id="example1">

```go {.bad #example2}
package util
```

Would result in a code element such as: <code class="language-go bad" id="example2">

These can be used to provide different styles for different languages, as well as
different styles (for example, background color) for code blocks. The IDs can be used
to select (or point to) specific code blocks in the rendered HTML.
*/
package ktw
