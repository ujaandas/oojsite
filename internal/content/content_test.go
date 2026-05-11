package content

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"oojsite/internal/templates"
)

func TestRenderPostsWithoutTemplate(t *testing.T) {
	root := t.TempDir()
	postsDir := filepath.Join(root, "posts")
	outDir := filepath.Join(root, "out")
	for _, dir := range []string{postsDir, outDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	postPath := filepath.Join(postsDir, "blog", "hello.md")
	writeFile(t, postPath, "---\ntitle: Hello\n---\n# Hello\n\nThis is markdown.")

	posts, err := LoadPosts(postsDir)
	if err != nil {
		t.Fatalf("LoadPosts: %v", err)
	}

	if err := RenderPosts(posts, outDir, template.New("")); err != nil {
		t.Fatalf("RenderPosts: %v", err)
	}

	outPath := filepath.Join(outDir, "posts", "blog", "hello", "index.html")
	output := readFile(t, outPath)
	if strings.Contains(output, "<!DOCTYPE html>") {
		t.Fatalf("expected raw markdown rendering, got layout wrapper: %s", output)
	}
	if !strings.Contains(output, "<h1>Hello</h1>") || !strings.Contains(output, "<p>This is markdown.</p>") {
		t.Fatalf("expected rendered markdown HTML, got: %s", output)
	}
}

func TestRenderPostsWithTemplate(t *testing.T) {
	root := t.TempDir()
	tmplDir := filepath.Join(root, "templates")
	componentDir := filepath.Join(root, "components")
	siteDir := filepath.Join(root, "site")
	postsDir := filepath.Join(root, "posts")
	outDir := filepath.Join(root, "out")
	for _, dir := range []string{tmplDir, componentDir, siteDir, postsDir, outDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	writeFile(t, filepath.Join(tmplDir, "custom.html"), `<html><body>{{ get .Frontmatter "title" }}::{{ .Content }}</body></html>`)
	writeFile(t, filepath.Join(postsDir, "note.md"), "---\ntitle: Custom Post\ntemplate: custom\n---\n# Body")

	tmpls, err := templates.Load(tmplDir, componentDir, siteDir)
	if err != nil {
		t.Fatalf("templates.Load: %v", err)
	}

	posts, err := LoadPosts(postsDir)
	if err != nil {
		t.Fatalf("LoadPosts: %v", err)
	}

	if err := RenderPosts(posts, outDir, tmpls); err != nil {
		t.Fatalf("RenderPosts: %v", err)
	}

	output := readFile(t, filepath.Join(outDir, "posts", "note", "index.html"))
	if !strings.Contains(output, "Custom Post::") {
		t.Fatalf("expected custom template output, got: %s", output)
	}
	if !strings.Contains(output, "<h1>Body</h1>") {
		t.Fatalf("expected rendered markdown in template content, got: %s", output)
	}
}

func TestRenderPostsWithoutFrontmatter(t *testing.T) {
	root := t.TempDir()
	postsDir := filepath.Join(root, "posts")
	outDir := filepath.Join(root, "out")
	for _, dir := range []string{postsDir, outDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	writeFile(t, filepath.Join(postsDir, "plain.md"), "# Plain\n\nNo frontmatter here.")

	posts, err := LoadPosts(postsDir)
	if err != nil {
		t.Fatalf("LoadPosts: %v", err)
	}

	if err := RenderPosts(posts, outDir, template.New("")); err != nil {
		t.Fatalf("RenderPosts: %v", err)
	}

	output := readFile(t, filepath.Join(outDir, "posts", "plain", "index.html"))
	if !strings.Contains(output, "<h1>Plain</h1>") {
		t.Fatalf("expected markdown output without frontmatter, got: %s", output)
	}
}

func TestLoadPostsRejectsMalformedFrontmatter(t *testing.T) {
	root := t.TempDir()
	postsDir := filepath.Join(root, "posts")
	if err := os.MkdirAll(postsDir, 0755); err != nil {
		t.Fatalf("mkdir %s: %v", postsDir, err)
	}

	writeFile(t, filepath.Join(postsDir, "broken.md"), "---\ntitle: missing end marker")

	if _, err := LoadPosts(postsDir); err == nil {
		t.Fatal("expected LoadPosts to fail for malformed frontmatter")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
