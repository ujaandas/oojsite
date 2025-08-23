package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

func CleanOutput(outDir string) error {
	if err := os.RemoveAll(outDir); err != nil {
		return fmt.Errorf("removing %q: %w", outDir, err)
	}
	return nil
}

func GenerateSite(
	posts []PostMeta,
	outDir string,
	tpl *template.Template,
	logger *log.Logger,
) error {
	for _, p := range posts {
		outPath := filepath.Join(outDir, p.Slug+".html")
		logger.Printf("rendering post → %s", outPath)
		if err := RenderPost(tpl, PostMeta(p), outPath); err != nil {
			return fmt.Errorf("render post %q: %w", p.Slug, err)
		}
	}

	idxPath := filepath.Join(outDir, "index.html")
	logger.Printf("rendering index → %s", idxPath)

	metas := make([]PostMeta, len(posts))
	for i, p := range posts {
		metas[i] = PostMeta(p)
	}
	if err := RenderIndex(tpl, metas, idxPath); err != nil {
		return fmt.Errorf("render index: %w", err)
	}

	return nil
}
