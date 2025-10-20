package main

import (
	"bufio"
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
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
)

var tagPostMap = make(map[string][]BlogPost) // tag -> posts

type Frontmatter struct {
	Title    string         `yaml:"title"`
	Template string         `yaml:"template"`
	Tags     []string       `yaml:"tags,omitempty"`
	Misc     map[string]any `yaml:",inline"`
}

type BlogPost struct {
	Meta    Frontmatter
	Snippet string
	Raw     []byte
}

type BlogTemplate struct {
	Title   string
	Content template.HTML
	Misc    map[string]any
}

type PageTemplate map[string][]BlogPost

func init() {
	flag.StringVar(&outFlag, "out", "out", "where to generate outputted site")
}

func main() {
	flag.Parse()

	// load templates
	tmpls, err := template.ParseFS(tmplFS, "templates/*.html", "site/*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	// create output directory
	if err := os.MkdirAll(outFlag, os.ModePerm); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	// copy public dir
	copyContents("public", fmt.Sprintf("%s/public", outFlag))

	// process markdown
	filepath.Walk("site", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".md") {
			processMarkdown(path, tmpls)
		}
		return nil
	})

	// process pages
	fs.WalkDir(pageFS, "site", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".html") {
			processHTMLPage(path, tmpls)
		}
		return nil
	})
}

func copyContents(srcDir, dstDir string) {
	err := filepath.Walk(srcDir, func(srcPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			log.Fatalf("walk error at %s: %v", srcPath, walkErr)
		}

		// derive relative path to mirror structure in dstDir
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			log.Fatalf("cannot compute relative path for %s: %v", srcPath, err)
		}
		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			// ensure directory exists at destination
			if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
				log.Fatalf("mkdir failed for %s: %v", dstPath, err)
			}
			return nil
		}

		// open source file for reading
		srcFile, err := os.Open(srcPath)
		if err != nil {
			log.Fatalf("cannot open source file %s: %v", srcPath, err)
		}
		defer srcFile.Close()

		// ensure parent directory exists for destination file
		if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
			log.Fatalf("mkdir failed for parent of %s: %v", dstPath, err)
		}

		// create destination file for writing
		dstFile, err := os.Create(dstPath)
		if err != nil {
			log.Fatalf("cannot create destination file %s: %v", dstPath, err)
		}
		defer dstFile.Close()

		// copy file contents from src to dst
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			log.Fatalf("copy failed from %s to %s: %v", srcPath, dstPath, err)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("directory walk failed: %v", err)
	}
}

func processHTMLPage(path string, pages *template.Template) {
	// get filename
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

func processMarkdown(path string, tmpls *template.Template) {
	// read file
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read markdown file %s: %v", path, err)
	}

	fileContent := extractFileContent(path, content)

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

	dirOutPath := filepath.Join(outFlag, filepath.Dir(pathWoSite))
	if err := os.MkdirAll(dirOutPath, os.ModePerm); err != nil {
		log.Fatalf("failed to create dirs %s: %v", filepath.Dir(pathWoSite), err)
	}

	outPath := filepath.Join(outFlag, strings.TrimSuffix(pathWoSite, ".md")+".html")
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
		Misc:    fileContent.Meta.Misc,
	})
	if err != nil {
		log.Fatalf("failed to execute template for %s: %v", path, err)
	}

	// collect tags and map posts
	for _, tag := range fileContent.Meta.Tags {
		tagPostMap[tag] = append(tagPostMap[tag], *fileContent)
	}
}

func extractFileContent(path string, content []byte) *BlogPost {
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

	// TODO: converts parts[2] to actual char, and implement snippet
	return &BlogPost{
		Meta:    meta,
		Snippet: makeSnippet(parts[2]),
		Raw:     parts[2],
	}
}

func makeSnippet(raw []byte) string {
	const wordCount = 20
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	words := []string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		for word := range strings.FieldsSeq(line) {
			words = append(words, word)
			if len(words) >= wordCount {
				break
			}
		}
		if len(words) >= wordCount {
			break
		}
	}

	return strings.Join(words, " ") + "..."
}

func sortedPosts(posts []BlogPost) []BlogPost {
	sorted := make([]BlogPost, len(posts))
	copy(sorted, posts)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Meta.Title < sorted[j].Meta.Title // sort by title
	})
	return sorted
}
