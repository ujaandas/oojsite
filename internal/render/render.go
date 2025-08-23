package render

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

const (
	IndexTpl = "index.html"
	PostTpl  = "post.html"
)

func RenderIndex(tpl *template.Template, posts []PostMeta, outPath string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create index %q: %w", outPath, err)
	}
	defer f.Close()

	data := struct {
		Posts []PostMeta
	}{
		Posts: posts,
	}

	if err := tpl.ExecuteTemplate(f, IndexTpl, data); err != nil {
		return fmt.Errorf("executing %s: %w", IndexTpl, err)
	}
	return nil
}

func RenderPost(tpl *template.Template, post PostMeta, outPath string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("mkdir for %q: %w", outPath, err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create post %q: %w", outPath, err)
	}
	defer f.Close()

	if err := tpl.ExecuteTemplate(f, PostTpl, post); err != nil {
		return fmt.Errorf("executing %s: %w", PostTpl, err)
	}
	return nil
}

type PostMeta struct {
	Slug    string
	Title   string
	Content string
}
