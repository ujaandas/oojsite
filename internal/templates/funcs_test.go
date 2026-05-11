package templates

import (
	"testing"

	"oojsite/internal/model"
)

func TestFuncsRegisterExpectedHelpers(t *testing.T) {
	funcs := Funcs()
	required := []string{
		"groupBy",
		"sortBy",
		"sortByDesc",
		"filter",
		"findPostByField",
		"filterPostsByField",
		"findPostByNotField",
		"filterPostsByNotField",
		"first",
		"reverse",
		"unique",
		"get",
		"slugify",
		"truncate",
		"formatDate",
	}

	for _, name := range required {
		if _, ok := funcs[name]; !ok {
			t.Fatalf("expected helper %q to be registered", name)
		}
	}
}

func TestFindPostAndExcludeHelpers(t *testing.T) {
	posts := []model.Post{
		{Frontmatter: map[string]interface{}{"slug": "home", "tags": []interface{}{"home"}, "title": "Home"}},
		{Frontmatter: map[string]interface{}{"slug": "guide", "tags": []interface{}{"nav-hidden"}, "title": "Guide"}},
		{Frontmatter: map[string]interface{}{"slug": "post", "tags": []interface{}{"docs"}, "title": "Post"}},
	}

	if got := findFirstPostByField("slug", "guide", posts); got == nil || getFieldValue(got.Frontmatter, "title") != "Guide" {
		t.Fatalf("expected to find guide by field")
	}

	if got := findFirstPostByNotField("slug", "home", posts); got == nil || getFieldValue(got.Frontmatter, "title") != "Guide" {
		t.Fatalf("expected to find first post not matching slug")
	}

	if got := findPostsByField("tags", "docs", posts); len(got) != 1 {
		t.Fatalf("expected one docs-tagged post, got %d", len(got))
	}
}

func TestCollectionAndStringHelpers(t *testing.T) {
	posts := []model.Post{
		{Frontmatter: map[string]interface{}{"tags": []interface{}{"docs", "go"}, "date": "January 2, 2024", "title": "Second"}},
		{Frontmatter: map[string]interface{}{"tags": []interface{}{"docs"}, "date": "January 1, 2024", "title": "First"}},
		{Frontmatter: map[string]interface{}{"tags": []interface{}{"blog"}, "date": "January 3, 2024", "title": "Third"}},
	}

	grouped := groupBy("tags", posts)
	if len(grouped["docs"]) != 2 {
		t.Fatalf("expected docs tag group to contain 2 posts, got %d", len(grouped["docs"]))
	}

	sorted := sortByDesc("date", posts)
	if getFieldValue(sorted[0].Frontmatter, "title") != "Third" {
		t.Fatalf("expected newest post first, got %q", getFieldValue(sorted[0].Frontmatter, "title"))
	}

	if got := unique("tags", posts); len(got) != 3 {
		t.Fatalf("expected 3 unique tags, got %d", len(got))
	}

	if got := slugify("My Blog Post!"); got != "my-blog-post" {
		t.Fatalf("unexpected slugify result: %q", got)
	}

	if got := truncate(4, "trim me"); got != "trim..." {
		t.Fatalf("unexpected truncate result: %q", got)
	}

	if got := formatDate("2006-01-02", "January 2, 2024"); got != "2024-01-02" {
		t.Fatalf("unexpected formatted date: %q", got)
	}

	if got := getSafe(map[string]interface{}{"title": "Hello"}, "title"); got != "Hello" {
		t.Fatalf("unexpected getSafe result: %q", got)
	}

	if got := getSafe(map[string]interface{}{}, "missing"); got != "" {
		t.Fatalf("expected empty string for missing key, got %q", got)
	}
}
