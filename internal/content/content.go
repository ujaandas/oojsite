package content

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kaleocheng/goldmark"
	"github.com/kaleocheng/goldmark/ast"
	"github.com/kaleocheng/goldmark/text"
	"gopkg.in/yaml.v2"

	"oojsite/internal/model"
)

func LoadPosts(postDir string) ([]model.Post, error) {
	var posts []model.Post

	err := filepath.Walk(postDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}

		post, err := loadPost(path, postDir)
		if err != nil {
			return err
		}

		posts = append(posts, *post)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].SourcePath < posts[j].SourcePath
	})

	return posts, nil
}

func RenderPosts(posts []model.Post, outDir string, tmpls *template.Template) error {
	for _, post := range posts {
		if err := renderPost(post, posts, outDir, tmpls); err != nil {
			return err
		}
	}
	return nil
}

func RenderPages(pageDir, outDir string, posts []model.Post, tmpls *template.Template) error {
	return filepath.Walk(pageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}

		rel, err := filepath.Rel(pageDir, path)
		if err != nil {
			return err
		}

		return renderPage(rel, outDir, posts, tmpls)
	})
}

func loadPost(path, postDir string) (*model.Post, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	post, err := extractFrontmatter(path, content)
	if err != nil {
		return nil, err
	}

	rel, err := filepath.Rel(postDir, path)
	if err != nil {
		return nil, err
	}
	post.SourcePath = path
	post.OutputRel = strings.TrimSuffix(rel, filepath.Ext(rel))
	post.Filepath = "/posts/" + filepath.ToSlash(post.OutputRel) + "/"
	return post, nil
}

func renderPost(post model.Post, posts []model.Post, outDir string, tmpls *template.Template) error {
	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert(post.Raw, &buf); err != nil {
		return err
	}

	templateName := "post.html"
	if userTemplate, ok := post.Frontmatter["template"]; ok {
		if templateStr, isString := userTemplate.(string); isString && templateStr != "" {
			templateName = templateStr
			if !strings.HasSuffix(templateName, ".html") {
				templateName += ".html"
			}
		}
	}

	outPath := filepath.Join(outDir, "posts", post.OutputRel, "index.html")
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	data := model.TemplateData{
		Frontmatter: post.Frontmatter,
		Content:     template.HTML(buf.String()),
		Global:      model.GlobalData{Posts: posts},
	}

	selected := tmpls.Lookup(templateName)
	if selected == nil {
		return fmt.Errorf("template %s not found", templateName)
	}
	return selected.Execute(outFile, data)
}

func renderPage(path, outDir string, posts []model.Post, tmpls *template.Template) error {
	outPath := filepath.Join(outDir, path)
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}

	tmpl := tmpls.Lookup(path)
	if tmpl == nil {
		return fmt.Errorf("template %s not found", path)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	data := model.PageData{Global: model.GlobalData{Posts: posts}}
	return tmpl.Execute(outFile, data)
}

func extractFrontmatter(path string, content []byte) (*model.Post, error) {
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("no frontmatter found in %s", path)
	}

	rawFrontmatter := parts[1]
	rawBody := parts[2]

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal(rawFrontmatter, &frontmatter); err != nil {
		return nil, fmt.Errorf("failed to parse front matter in %s: %v", path, err)
	}
	if frontmatter == nil {
		frontmatter = make(map[string]interface{})
	}

	return &model.Post{
		Frontmatter: frontmatter,
		Snippet:     makeSnippet(rawBody, 20),
		Raw:         rawBody,
	}, nil
}

func makeSnippet(raw []byte, wordCount int) string {
	text := extractText(raw)
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	if len(words) <= wordCount {
		return strings.Join(words, " ")
	}
	return strings.Join(words[:wordCount], " ") + "..."
}

func extractText(raw []byte) string {
	doc := goldmark.DefaultParser().Parse(text.NewReader(raw))
	var b strings.Builder

	var walk func(ast.Node)
	walk = func(n ast.Node) {
		if t, ok := n.(*ast.Text); ok {
			b.Write(t.Segment.Value(raw))
			b.WriteByte(' ')
		}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			walk(c)
		}
	}

	walk(doc)
	return b.String()
}
