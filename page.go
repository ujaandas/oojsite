package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

/*
Handle everything pertaining to our final outputted HTML, including post templates.
Confusingly, a page is _not_ simply just landing pages like `index.html`, your blog post
will also eventually become a "page" - basically anything that ends up as HTML is a page.
Additionally, all pages have access to Go's templating features, allowing end-users to
access a wide breadth of global state, things like tags, page titles, options, etc...
*/

var tagPostMap = make(map[string][]Post) // tag -> posts

type Template struct {
	Title   string
	Content template.HTML
}

type PageTemplate map[string][]Post

// Load both page templates (ie; for posts) and actual pages (ie; index.html).
func loadPages(tmplDir, siteDir string) (*template.Template, error) {
	tmpls := template.New("")

	// Load post templates
	tmpls, err := tmpls.ParseGlob(fmt.Sprintf("%s/*.html", tmplDir))
	if err != nil {
		return nil, err
	}

	// Load actual page templates
	tmpls, err = tmpls.ParseGlob(fmt.Sprintf("%s/*.html", siteDir))
	if err != nil {
		return nil, err
	}

	return tmpls, nil
}

func processPage(path, outDir string, pages *template.Template) error {
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

	return nil
}
