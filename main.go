package main

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Post struct {
	Slug    string
	Title   string
	Content string
}

func main() {
	var posts []Post
	err := filepath.WalkDir("content/blog", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		slug := strings.TrimSuffix(d.Name(), ".md")
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(raw)
		title := extractTitle(content, slug)
		posts = append(posts, Post{Slug: slug, Title: title, Content: content})
		return nil
	})
	if err != nil {
		log.Fatalf("walking content/blog: %v", err)
	}

	if err := os.RemoveAll("public"); err != nil {
		log.Fatalf("remove public: %v", err)
	}
	if err := os.MkdirAll("public/blog", 0755); err != nil {
		log.Fatalf("mkdir public/blog: %v", err)
	}

	tpl := template.Must(template.ParseGlob("templates/*.html"))

	for _, p := range posts {
		outPath := filepath.Join("public/blog", p.Slug+".html")
		f, err := os.Create(outPath)
		if err != nil {
			log.Fatalf("create %s: %v", outPath, err)
		}
		if err := tpl.ExecuteTemplate(f, "post.html", p); err != nil {
			f.Close()
			log.Fatalf("execute post.html: %v", err)
		}
		f.Close()
		log.Printf("Rendered post: %s", outPath)
	}

	idx, err := os.Create("public/index.html")
	if err != nil {
		log.Fatalf("create index.html: %v", err)
	}
	data := struct{ Posts []Post }{posts}
	if err := tpl.ExecuteTemplate(idx, "index.html", data); err != nil {
		idx.Close()
		log.Fatalf("execute index.html: %v", err)
	}
	idx.Close()
	log.Println("Site generated at ./public")
}

func extractTitle(content, slug string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[2:])
		}
	}
	return slug
}
