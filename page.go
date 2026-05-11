package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
Handle everything pertaining to our final outputted HTML, including post templates.
Confusingly, a page is _not_ simply just landing pages like `index.html`, your blog post
will also eventually become a "page" - basically anything that ends up as HTML is a page.
Additionally, all pages have access to Go's templating features, allowing end-users to
access a wide breadth of global state through .Global.
*/

var allPosts []Post // Collects all posts as they're processed

type GlobalData struct {
	Posts []Post
}

type Template struct {
	Content     template.HTML
	Frontmatter map[string]interface{}
	Global      GlobalData
}

type PageData struct {
	Global GlobalData
}

// Load both page templates (ie; for posts) and actual pages (ie; index.html).
func loadTemplates(tmplDir, componentDir, siteDir string) (*template.Template, error) {
	tmpls := template.New("")

	// Register custom template functions
	tmpls.Funcs(createTemplateFuncs())

	// load post templates
	filepath.Walk(tmplDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		rel, err := filepath.Rel(tmplDir, path)
		if err != nil {
			return err
		}

		content, _ := os.ReadFile(path)
		tmpls.New(rel).Parse(string(content))
		return err
	})

	// load actual page templates
	filepath.Walk(siteDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		rel, err := filepath.Rel(siteDir, path)
		if err != nil {
			return err
		}

		content, _ := os.ReadFile(path)
		tmpls.New(rel).Parse(string(content))
		return err
	})

	// load component templates
	filepath.Walk(componentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		rel, err := filepath.Rel(componentDir, path)
		if err != nil {
			return err
		}

		content, _ := os.ReadFile(path)
		tmpls.New(rel).Parse(string(content))
		return err
	})

	return tmpls, nil
}

func processPage(path, outDir string, tmpls *template.Template) error {
	// get expected output path
	outPath := filepath.Join(outDir, path)

	// ensure path exists
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}

	// apply template
	tmpl := tmpls.Lookup(path)
	if tmpl == nil {
		log.Fatalf("template %s not found for %s", path, path)
	}

	// create output file
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("failed to create output file %s: %v", outPath, err)
	}
	defer outFile.Close()

	data := PageData{
		Global: GlobalData{
			Posts: allPosts,
		},
	}

	// write output file
	err = tmpl.Execute(outFile, data)
	if err != nil {
		log.Fatalf("failed to execute template for %s: %v", path, err)
	}

	return nil
}
