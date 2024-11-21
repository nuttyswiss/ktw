package main

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/nuttyswiss/ktw"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func init() {
	var cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate the website",
		RunE:  generate,
	}
	cli.AddCommand(cmd)
}

// splitContent split the full contents of a file into the frontmatter
// delimiter, the metadata, and the rest of the content. It returns an
// error if delimiter is not found.
func splitContent(buf []byte) ([]byte, []byte, []byte, error) {
	delim, rest, present := bytes.Cut(buf, []byte("\n"))
	if !present {
		return nil, nil, nil, fmt.Errorf("delimiter not found")
	}
	delim = append(delim, '\n')

	meta, content, ok := bytes.Cut(rest, delim)
	if !ok {
		return nil, nil, nil, fmt.Errorf("delimiter not found")
	}

	return delim, meta, content, nil
}

// generate traverses a directory of files representing a web site. For each
// file that we encounter, if it is a file that we need to process, we go and
// process that file (usually generate an HTML file from Markdown).
func generate(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("generate takes no arguments")
	}
	root := viper.GetString("dir")
	if root == "" {
		return fmt.Errorf("config is missing 'dir' key")
	}

	fmt.Printf("Generating from %s\n", root)
	err := filepath.WalkDir(root, func(src string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		srcpath, err := filepath.Rel(root, src)
		if err != nil {
			return err
		}
		name := filepath.Base(srcpath)
		ext := filepath.Ext(name)
		if ext != ".md" {
			return nil
		}

		dstpath := srcpath[:len(srcpath)-len(ext)] + ".html"
		dst := filepath.Join(root, dstpath)
		fmt.Printf("Generate HTML: %s --> %s", srcpath, dstpath)

		inbuf, err := os.ReadFile(src)
		if err != nil {
			return err
		}

		delim, frontmatter, content, err := splitContent(inbuf)
		if err != nil {
			return fmt.Errorf("failed to split content in %q: %w", srcpath, err)
		}
		metadata := make(map[string]any)
		if bytes.HasPrefix(delim, []byte("---")) {
			if err := yaml.Unmarshal(frontmatter, metadata); err != nil {
				return fmt.Errorf("failed to parse metadata in %q: %w", srcpath, err)
			}
		}

		page := ktw.Page{
			Title:    name,
			Metadata: metadata,
			Contents: []ktw.Renderer{ktw.Markdown(content)},
		}
		var outbuf bytes.Buffer
		if err := page.Render(context.Background(), &outbuf); err != nil {
			return err
		}
		fmt.Println(", Done!")

		if err := os.WriteFile(dst, outbuf.Bytes(), 0644); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
