package render

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

//go:embed templates/*html
var defaultTemplates embed.FS

func LoadTemplates(tplDir string) (*template.Template, error) {
	root := template.New("")
	root, err := root.ParseFS(defaultTemplates, "templates/*.html")
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
