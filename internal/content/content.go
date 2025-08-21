package content

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type PostMeta struct {
	Slug    string
	Title   string
	Content string
}

func DiscoverPosts(srcDir string) ([]PostMeta, error) {
	var posts []PostMeta

	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking %q: %w", path, err)
		}
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %q: %w", path, err)
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("rel path: %w", err)
		}
		slug := strings.TrimSuffix(rel, filepath.Ext(rel))
		slug = filepath.ToSlash(slug)

		title := extractTitle(string(raw), slug)

		posts = append(posts, PostMeta{
			Slug:    slug,
			Title:   title,
			Content: string(raw),
		})
		return nil
	})

	if err != nil {
		return nil, err
	}
	return posts, nil
}

func extractTitle(content, slug string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[2:])
		}
	}
	return slug
}
