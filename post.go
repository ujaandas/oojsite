package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kaleocheng/goldmark"
	"github.com/kaleocheng/goldmark/ast"
	"github.com/kaleocheng/goldmark/text"
	"gopkg.in/yaml.v2"
)

/*
A "post", in the context of this app, refers almost exclusively to markdown (`*.md`) files.
For a given markdown file, we need to handle 4 main things:
	1. Frontmatter processing
	2. Markdown to HTML conversion
	3. Apply template
	4. Write the output
*/

type Frontmatter struct {
	Title    string   `yaml:"title"`
	Template string   `yaml:"template"`
	Tags     []string `yaml:"tags,omitempty"`
}

type Post struct {
	Meta     Frontmatter
	Filepath string
	Snippet  string
	Raw      []byte
}

// Process and extract frontmatter, converting markdown file to our struct.
func extractFrontmatter(path string, content []byte) (*BlogPost, error) {
	// Split frontmatter
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("no frontmatter found in %s", content)
	}

	// Split parts (parts[0] is preamble, should be empty and we don't care about it anyways)
	frontmatter := parts[1]
	body := parts[2]

	// Unmarshal and read
	var meta Frontmatter
	if err := yaml.Unmarshal(frontmatter, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse front matter in %s: %v", path, err)
	}

	// Get relative filepath
	relFp := strings.TrimPrefix(strings.TrimSuffix(path, filepath.Ext(path))+".html", "site/")

	return &BlogPost{
		Meta:     meta,
		Filepath: relFp,
		Snippet:  makeSnippet(body, 20),
		Raw:      body,
	}, nil
}

// Take the raw text and form a snippet of `wordCount` words.
func makeSnippet(raw []byte, wordCount int) string {
	text := extractText(raw)
	words := strings.Fields(text)

	if len(words) <= wordCount {
		return strings.Join(words, " ")
	}
	return strings.Join(words[:wordCount], " ") + "..."
}

// Walk through markdown AST and only extract text nodes.
func extractText(raw []byte) string {
	doc := goldmark.DefaultParser().Parse(text.NewReader(raw))

	// Hold words here
	var b strings.Builder

	// Walk through AST with DFS
	var walk func(ast.Node)
	walk = func(n ast.Node) {
		if t, ok := n.(*ast.Text); ok {
			b.Write(t.Segment.Value(raw))
			b.WriteByte(' ') // Ensure spaces between words
		}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			walk(c)
		}
	}

	walk(doc)
	return b.String()
}
