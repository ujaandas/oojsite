package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kaleocheng/goldmark"
)

//go:embed templates/*.html site/*.html
var tmplFS embed.FS

//go:embed site/*.html
var pageFS embed.FS

var tagPostMap = make(map[string][]BlogPost) // tag -> posts

type BlogPost struct {
	Meta     Frontmatter
	Filepath string
	Snippet  string
	Raw      []byte
}

type BlogTemplate struct {
	Title   string
	Content template.HTML
	Misc    map[string]any
}

type PageTemplate map[string][]BlogPost

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
			processMarkdown(path, cfg.outDir, tmpls)
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

func processMarkdown(path, outDir string, tmpls *template.Template) {
	// read file
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read markdown file %s: %v", path, err)
	}

	fileContent, err := extractFrontmatter(path, content)
	if err != nil {
		log.Fatalf("Failed to extract frontmatter: %v", err)
	}

	// convert markdown to HTML
	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert(fileContent.Raw, &buf); err != nil {
		log.Fatalf("failed to convert markdown in %s: %v", path, err)
	}

	// apply template
	tmpl := tmpls.Lookup(fmt.Sprintf("%s.html", fileContent.Meta.Template))
	if tmpl == nil {
		log.Fatalf("template %s not found for %s", fileContent.Meta.Template, path)
	}

	// create output file + parent dirs
	trimmedPath := strings.TrimPrefix(path, "site/")
	pathWoSite := filepath.ToSlash(trimmedPath)

	dirOutPath := filepath.Join(outDir, filepath.Dir(pathWoSite))
	if err := os.MkdirAll(dirOutPath, os.ModePerm); err != nil {
		log.Fatalf("failed to create dirs %s: %v", filepath.Dir(pathWoSite), err)
	}

	outPath := filepath.Join(outDir, strings.TrimSuffix(pathWoSite, ".md")+".html")
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("failed to create output file %s: %v", outPath, err)
	}
	defer outFile.Close()

	// write output file
	// TODO: automatch (by name) everything in fileContent.meta
	err = tmpl.Execute(outFile, BlogTemplate{
		Title:   fileContent.Meta.Title,
		Content: template.HTML(buf.String()),
		// Misc:    fileContent.Meta.Misc,
	})
	if err != nil {
		log.Fatalf("failed to execute template for %s: %v", path, err)
	}

	// collect tags and map posts
	for _, tag := range fileContent.Meta.Tags {
		tagPostMap[tag] = append(tagPostMap[tag], *fileContent)
	}
}

func sortedPosts(posts []BlogPost) []BlogPost {
	sorted := make([]BlogPost, len(posts))
	copy(sorted, posts)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Meta.Title < sorted[j].Meta.Title // sort by title
	})
	return sorted
}
