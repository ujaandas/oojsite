package main

import (
	"fmt"
	"log"
	"net/http"
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
	tmpls, err := loadTemplates(cfg.templateDir, cfg.componentDir, cfg.pageDir)
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}
	log.Println("Templates loaded!")

	// Process markdown
	filepath.Walk(cfg.postDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".md") {
			log.Printf("Processing post at %s...\n", path)

			rel, err := filepath.Rel(cfg.postDir, path)
			if err != nil {
				return err
			}

			// Pass both
			if err := processPost(rel, cfg.postDir, cfg.outDir, tmpls); err != nil {
				log.Fatalf("Failed to process post %s: %v", path, err)
			}
		}
		return nil
	})

	// Process pages
	filepath.Walk(cfg.pageDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			log.Printf("Processing page at %s...\n", path)

			rel, err := filepath.Rel(cfg.pageDir, path)
			if err != nil {
				return err
			}

			if err := processPage(rel, cfg.outDir, tmpls); err != nil {
				log.Fatalf("Failed to process page %s: %v", path, err)
			}

			log.Println("Page processed!")
		}
		return nil
	})

	// Build and compile TailwindCSS
	log.Println("Building TailwindCSS...")
	if err := buildTailwind(cfg.outDir, cfg.staticDir); err != nil {
		log.Fatalf("Failed to build TailwindCSS: %v", err)
	}
	log.Println("TailwindCSS built!")

	// Copy over static files
	log.Println("Copying static files...")
	if err := copyStaticContents(cfg.staticDir, fmt.Sprintf("%s/static", cfg.outDir)); err != nil {
		log.Fatalf("Failed to copy static files: %v", err)
	}
	log.Println("Copied static files!")

	// Serve
	if !cfg.dev {
		return
	}

	log.Println("Server started on localhost:8000!")
	if err := http.ListenAndServe(":8000", http.FileServer(http.Dir(cfg.outDir))); err != nil {
		log.Fatalf("Server has crashed: %v", err)
	}
}
