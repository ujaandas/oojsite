package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
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
	Frontmatter Frontmatter
	Filepath    string
	Snippet     string
	Raw         []byte
}

func processPost(path, outDir string, tmpls *template.Template) error {
	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Extract frontmatter
	post, err := extractFrontmatter(path, content)
	if err != nil {
		return err
	}

	// Convert markdown to HTML
	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert(post.Raw, &buf); err != nil {
		return err
	}

	// Apply template
	tmpl := tmpls.Lookup(fmt.Sprintf("%s.html", post.Frontmatter.Template))
	if tmpl == nil {
		return err
	}

	// Write outputted file
	if err := writePostFile(path, outDir, tmpl, post, buf); err != nil {
		return err
	}

	return nil
}

func sortedPosts(posts []Post) []Post {
	sorted := make([]Post, len(posts))
	copy(sorted, posts)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Frontmatter.Title < sorted[j].Frontmatter.Title // sort by title
	})
	return sorted
}

// Process and extract frontmatter, converting markdown file to our struct.
func extractFrontmatter(path string, content []byte) (*Post, error) {
	// Split frontmatter
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("no frontmatter found in %s", content)
	}

	// Split parts (parts[0] is preamble, should be empty and we don't care about it anyways)
	rawFrontmatter := parts[1]
	rawBody := parts[2]

	// Unmarshal and read
	var frontmatter Frontmatter
	if err := yaml.Unmarshal(rawFrontmatter, &frontmatter); err != nil {
		return nil, fmt.Errorf("failed to parse front matter in %s: %v", path, err)
	}

	// Get relative filepath
	relFp := strings.TrimPrefix(strings.TrimSuffix(path, filepath.Ext(path))+".html", "site/")

	return &Post{
		Frontmatter: frontmatter,
		Filepath:    relFp,
		Snippet:     makeSnippet(rawBody, 20),
		Raw:         rawBody,
	}, nil
}

// Take the raw text and form a snippet of `wordCount` words.
func makeSnippet(raw []byte, wordCount int) string {
	text := extractText(raw)
	words := strings.Fields(text)

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

// Write outputted HTML file
func writePostFile(src, dst string, tmpl *template.Template, post *Post, contentBuf bytes.Buffer) error {
	// Create output file + parent dirs
	trimmedPath := strings.TrimPrefix(src, "site/")
	pathWoSite := filepath.ToSlash(trimmedPath)

	dirOutPath := filepath.Join(dst, filepath.Dir(pathWoSite))
	if err := os.MkdirAll(dirOutPath, os.ModePerm); err != nil {
		return err
	}

	outPath := filepath.Join(dst, strings.TrimSuffix(pathWoSite, ".md")+".html")
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Write output file
	if err := tmpl.Execute(outFile, BlogTemplate{
		Title:   post.Frontmatter.Title,
		Content: template.HTML(contentBuf.String()),
	}); err != nil {
		return err
	}

	// collect tags and map posts
	for _, tag := range post.Frontmatter.Tags {
		tagPostMap[tag] = append(tagPostMap[tag], *post)
	}

	return nil
}
