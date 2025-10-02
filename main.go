package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaleocheng/goldmark"
	"gopkg.in/yaml.v2"
)

//go:embed templates/*.html site/*.html
var tmplFS embed.FS

//go:embed site/*.html
var pageFS embed.FS

var (
	outFlag string
	siteUrl string
)

type Frontmatter struct {
	Title    string `yaml:"title"`
	Template string `yaml:"template"`
}

type FileContent struct {
	meta    Frontmatter
	content []byte
}

func init() {
	flag.StringVar(&outFlag, "out", "out", "where to generate outputted site")
}

func main() {
	flag.Parse()
	fmt.Printf("out flag: %s\n", outFlag)

	// load templates
	tmpls, err := template.ParseFS(tmplFS, "templates/*.html", "site/*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	// create output directory
	if err := os.MkdirAll(outFlag, os.ModePerm); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	// process pages
	fs.WalkDir(pageFS, "site", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".html") {
			processHTMLPage(path, tmpls)
		}
		return nil
	})

	// process markdown
	filepath.Walk("site", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".md") {
			processMarkdown(path, tmpls)
		}
		return nil
	})
}

type indexTmpl struct {
	Posts []Frontmatter
}

func processHTMLPage(path string, pages *template.Template) {
	filename := filepath.Base(path)

	// apply template
	tmpl := pages.Lookup(filename)
	if tmpl == nil {
		log.Fatalf("template %s not found for %s", filename, path)
	}

	// create output file
	outPath := filepath.Join(outFlag, filename)
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("failed to create output file %s: %v", outPath, err)
	}
	defer outFile.Close()

	// write output file
	err = tmpl.Execute(outFile, indexTmpl{
		[]Frontmatter{
			{Title: "test"},
			{Title: "test2"},
		},
	})
	if err != nil {
		log.Printf("failed to execute template for %s: %v", path, err)
	}
}

func processFileContent(path string, content []byte) *FileContent {
	// split frontmatter
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		log.Fatalf("no frontmatter found in %s", content)
	}

	// unmarshal and read
	var meta Frontmatter
	if err := yaml.Unmarshal(parts[1], &meta); err != nil {
		log.Fatalf("failed to parse front matter in %s: %v", path, err)
	}

	return &FileContent{
		meta:    meta,
		content: parts[2],
	}
}

func processMarkdown(path string, tmpls *template.Template) {
	// read file
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read markdown file %s: %v", path, err)
	}

	fileContent := processFileContent(path, content)

	// convert markdown to HTML
	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert(fileContent.content, &buf); err != nil {
		log.Fatalf("failed to convert markdown in %s: %v", path, err)
	}

	// apply template
	tmpl := tmpls.Lookup(fmt.Sprintf("%s.html", fileContent.meta.Template))
	if tmpl == nil {
		log.Fatalf("template %s not found for %s", fileContent.meta.Template, path)
	}

	// create output file + parent dirs
	trimmedPath := strings.TrimPrefix(path, "site/")
	pathWoSite := filepath.ToSlash(trimmedPath)

	dirOutPath := filepath.Join(outFlag, filepath.Dir(pathWoSite))
	if err := os.MkdirAll(dirOutPath, os.ModePerm); err != nil {
		log.Printf("failed to create dirs %s: %v", filepath.Dir(pathWoSite), err)
	}

	outPath := filepath.Join(outFlag, strings.TrimSuffix(pathWoSite, ".md")+".html")
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("failed to create output file %s: %v", outPath, err)
	}
	defer outFile.Close()

	// write output file
	err = tmpl.Execute(outFile, map[string]any{
		"Title":   fileContent.meta.Title,
		"Content": template.HTML(buf.String()),
	})
	if err != nil {
		log.Fatalf("failed to execute template for %s: %v", path, err)
	}
}
