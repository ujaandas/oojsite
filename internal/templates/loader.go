package templates

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

func Load(tmplDir, componentDir, siteDir string) (*template.Template, error) {
	tmpls := template.New("").Funcs(Funcs())

	if err := parseDir(tmpls, tmplDir); err != nil {
		return nil, err
	}
	if err := parseDir(tmpls, siteDir); err != nil {
		return nil, err
	}
	if err := parseDir(tmpls, componentDir); err != nil {
		return nil, err
	}

	return tmpls, nil
}

func parseDir(tmpls *template.Template, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = tmpls.New(rel).Parse(string(content))
		return err
	})
}
