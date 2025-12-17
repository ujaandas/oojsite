package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*.html site/*.html
var tmplFS embed.FS

//go:embed site/*.html
var pageFS embed.FS

var tagPostMap = make(map[string][]Post) // tag -> posts

type Template struct {
	Title   string
	Content template.HTML
}

type PageTemplate map[string][]Post

func main() {
	cfg, err := parseFlags()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// load templates
	tmpls, err := template.ParseFS(tmplFS, "templates/*.html", "site/*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	// Build and compile TailwindCSS
	if err := buildTailwind(cfg.outDir); err != nil {
		log.Fatalf("Failed to build TailwindCSS: %v", err)
	}

	// Copy over static files
	if err := copyStaticContents("static", fmt.Sprintf("%s/static", cfg.outDir)); err != nil {
		log.Fatalf("Failed to copy static files: %v", err)
	}

	// process markdown
	filepath.Walk("site", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".md") {
			if err := processPost(path, fmt.Sprintf("%s/posts", cfg.outDir), tmpls); err != nil {
				log.Fatalf("Failed to process markdown file %s: %v", path, err)
			}
		}
		return nil
	})

	// process pages
	fs.WalkDir(pageFS, "site", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".html") {
			processHTMLPage(path, cfg.outDir, tmpls)
		}
		return nil
	})
}

func processHTMLPage(path, outDir string, pages *template.Template) {
	// get filename
	filename := filepath.Base(path)

	// apply template
	tmpl := pages.Lookup(filename)
	if tmpl == nil {
		log.Fatalf("template %s not found for %s", filename, path)
	}

	// create output file
	outPath := filepath.Join(outDir, filename)
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("failed to create output file %s: %v", outPath, err)
	}
	defer outFile.Close()

	// fill in tags
	data := make(PageTemplate)
	for tag, posts := range tagPostMap {
		data[tag] = sortedPosts(posts)
	}

	// write output file
	err = tmpl.Execute(outFile, data)
	if err != nil {
		log.Fatalf("failed to execute template for %s: %v", path, err)
	}
}
