package assets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildTailwindSkipsMissingStylesheet(t *testing.T) {
	root := t.TempDir()
	staticDir := filepath.Join(root, "static")
	outDir := filepath.Join(root, "out")
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		t.Fatalf("mkdir %s: %v", staticDir, err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("mkdir %s: %v", outDir, err)
	}

	if err := BuildTailwind(outDir, staticDir); err != nil {
		t.Fatalf("expected nil error when styles.css is missing, got %v", err)
	}
}
