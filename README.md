# ktw/web -- Markdown Renderer and Publisher

Warning: this is not production ready, and **will** change.

This is my feable attempt to create my own Markdown publishing framework. It is
modeled after a number of projects that I've used and/or been involved with
over the past few years. I take inspiration freely from [Hugo] and [Caddy], as
well as various Markdown interpreters. It is not meant to be a solution that
fits anyone but myself, hence it'll be rather opinionated in what it implements
and supports. IE: almost nothing. :smile:

To install the `web` CLI:

```bash
$ go install github.com/nuttyswiss/ktw/cmd/web@latest
```

To run this, you point the CLI at site and do:

```bash
$ web --site ~/some/site verify
$ web --site ~/some/site generate
$ web --site ~/some/site publish
```

Where the site directory looks something like:

```
├── article.tmpl
├── config.yaml
└── htdocs
    ├── 404.html
    ├── 500.html
    ├── blog
    │   ├── article-01
    │   │   ├── index.html
    │   │   └── index.md
    │   ├── article-02
    │   │   ├── index.html
    │   │   └── index.md
    │   └── ...
    ├── index.html
    └── index.md
```

The `generate` sub-command will parse all the `.md` Markdown files and generate
the analogous `.html` HTML file alongside the Markdown file. The `verify`
command checks if there are any out of date files. And the `publish` command
will copy the generated website to its final destination.

The `config.yaml` file contains the following:

```yaml
# Config for github.com/nuttyswiss/ktw/cmd/web CLI
site: "example.com"
root: "sftp://host.example.com/tmp/testdir"
dir: "htdocs"
```

Note: the key names in this YAML file will likely change, as may the structure
of the file. At the current time, `web` takes a `--config` argument, which can
be used to point to the different config file. This can be used to publish the
web site to a different place (for canary deployments, etc):

```bash
$ web --config canary.yaml --site ~/some/site publish
```

## Future Work

- [ ] Finish up transition to using [html/template] and allow for use of
template variables within HTML templates as well as Markdown content.
- [ ] Integrate D2 parsing, maybe something like [github.com/FurqanSoftware/goldmark-d2],
or write our own to directly use [oss.terrastruct.com/d2].
- [ ] Write a code block parser that can use a "before"/"after" method to show
the diff of a changed code block. Possibly using some separator, say `:::`, or
something similar to separate before and after. Then using both colour and
other visual indicators (strike-through, bold, etc) to show where and what the
change in the code is.
- [ ] Write our own Markdown parser, or use [rsc.io/markdown] as a starting
point to write a smaller and more targeted Markdown parser and HTML converter.
- [ ] Write an index generator, such that there is an easy method to generate
an index of a set of pages.
- [ ] Write support for "default frontmatter"/"metadata", such that we do not
need to repeat ourselves ad'nauseum.

[Hugo]: https://gohugo.io/
[Caddy]: https://caddyserver.com/
[github.com/FurqanSoftware/goldmark-d2]: https://pkg.go.dev/github.com/FurqanSoftware/goldmark-d2
[oss.terrastruct.com/d2]: https://pkg.go.dev/oss.terrastruct.com/d2
[rsc.io/markdown]: https://pkg.go.dev/rsc.io/markdown
[html/template]: https://pkg.go.dev/html/template
