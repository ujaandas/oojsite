package main

import (
	"fmt"
	"html/template"
	"oojsite/assets"
	"os"
	"path/filepath"
)

func LoadTemplates(tplDir string) (*template.Template, error) {
	root, err := assets.Templates()
	if err != nil {
		return nil, fmt.Errorf("parsing embedded templates: %w", err)
	}

	if stat, err := os.Stat(tplDir); err == nil && stat.IsDir() {
		pattern := filepath.Join(tplDir, "*.html")
		if _, err := root.ParseGlob(pattern); err != nil {
			return nil, fmt.Errorf("parsing disk templates %q: %w", pattern, err)
		}
	}

	return root, nil
}
