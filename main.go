package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.Println("Parsing options...")
	cfg, err := parseOptions()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	log.Println("Options parsed!")

	// Load templates
	log.Println("Loading templates...")
	tmpls, err := loadPages(cfg.templateDir, cfg.pageDir)
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}
	log.Println("Templates loaded!")

	// Process markdown
	filepath.Walk(cfg.postDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".md") {
			log.Printf("Processing post at %s...\n", path)
			if err := processPost(path, fmt.Sprintf("%s/posts", cfg.outDir), tmpls); err != nil {
				log.Fatalf("Failed to process markdown file %s: %v", path, err)
			}
			log.Println("Post processed!")
		}
		return nil
	})

	// Process pages
	filepath.Walk(cfg.pageDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			log.Printf("Processing page at %s...\n", path)
			if err := processPage(path, cfg.outDir, tmpls); err != nil {
				log.Fatalf("Failed to process page %s: %v", path, err)
			}
			log.Println("Page processed!")
		}
		return nil
	})

	// Build and compile TailwindCSS
	log.Println("Building TailwindCSS...")
	if err := buildTailwind(cfg.outDir); err != nil {
		log.Fatalf("Failed to build TailwindCSS: %v", err)
	}
	log.Println("TailwindCSS built!")

	// Copy over static files
	log.Println("Copying static files...")
	if err := copyStaticContents("static", fmt.Sprintf("%s/static", cfg.outDir)); err != nil {
		log.Fatalf("Failed to copy static files: %v", err)
	}
	log.Println("Copied static files!")
}
