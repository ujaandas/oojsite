package templates

import (
	"fmt"
	"html/template"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"oojsite/internal/model"
)

func Funcs() template.FuncMap {
	return template.FuncMap{
		"groupBy":               groupBy,
		"sortBy":                sortBy,
		"sortByDesc":            sortByDesc,
		"filter":                filterPosts,
		"findPostByField":       findFirstPostByField,
		"filterPostsByField":    findPostsByField,
		"findPostByNotField":    findFirstPostByNotField,
		"filterPostsByNotField": findPostsByNotField,
		"first":                 first,
		"reverse":               reversePosts,
		"unique":                unique,
		"get":                   getSafe,
		"slugify":               slugify,
		"truncate":              truncate,
		"formatDate":            formatDate,
	}
}

func groupBy(field string, posts []model.Post) map[string][]model.Post {
	result := make(map[string][]model.Post)
	for _, post := range posts {
		if slice, ok := post.Frontmatter[field].([]interface{}); ok {
			for _, item := range slice {
				if str, isStr := item.(string); isStr {
					result[str] = append(result[str], post)
				}
			}
			continue
		}

		val := getFieldValue(post.Frontmatter, field)
		if val != "" {
			result[val] = append(result[val], post)
		}
	}
	return result
}

func sortBy(field string, posts []model.Post) []model.Post {
	sorted := make([]model.Post, len(posts))
	copy(sorted, posts)

	sort.SliceStable(sorted, func(i, j int) bool {
		numI, okI := getNumericFieldValue(sorted[i].Frontmatter, field)
		numJ, okJ := getNumericFieldValue(sorted[j].Frontmatter, field)
		if okI && okJ {
			return numI < numJ
		}

		valI := getFieldValue(sorted[i].Frontmatter, field)
		valJ := getFieldValue(sorted[j].Frontmatter, field)

		dateI, errI := time.Parse("January 2, 2006", valI)
		dateJ, errJ := time.Parse("January 2, 2006", valJ)
		if errI == nil && errJ == nil {
			return dateI.Before(dateJ)
		}

		return valI < valJ
	})

	return sorted
}

func sortByDesc(field string, posts []model.Post) []model.Post {
	sorted := sortBy(field, posts)
	reversePosts(sorted)
	return sorted
}

func filterPosts(field, value string, posts []model.Post) []model.Post {
	var result []model.Post
	for _, post := range posts {
		if matchesField(post, field, value) {
			result = append(result, post)
		}
	}
	return result
}

func findPostsByNotField(field, value string, posts []model.Post) []model.Post {
	var result []model.Post
	for _, post := range posts {
		if !matchesField(post, field, value) {
			result = append(result, post)
		}
	}
	return result
}

func findFirstPostByNotField(field, value string, posts []model.Post) *model.Post {
	for i := range posts {
		if !matchesField(posts[i], field, value) {
			return &posts[i]
		}
	}
	return nil
}

func findPostsByField(field, value string, posts []model.Post) []model.Post {
	var result []model.Post
	for _, post := range posts {
		if matchesField(post, field, value) {
			result = append(result, post)
		}
	}
	return result
}

func findFirstPostByField(field, value string, posts []model.Post) *model.Post {
	for i := range posts {
		if matchesField(posts[i], field, value) {
			return &posts[i]
		}
	}
	return nil
}

func first(n int, posts []model.Post) []model.Post {
	if n < 0 || n > len(posts) {
		return posts
	}
	return posts[:n]
}

func reversePosts(posts []model.Post) []model.Post {
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}
	return posts
}

func unique(field string, posts []model.Post) []string {
	seen := make(map[string]bool)
	var result []string

	for _, post := range posts {
		if slice, ok := post.Frontmatter[field].([]interface{}); ok {
			for _, item := range slice {
				if str, isStr := item.(string); isStr && !seen[str] {
					seen[str] = true
					result = append(result, str)
				}
			}
			continue
		}

		val := getFieldValue(post.Frontmatter, field)
		if val != "" && !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}

	return result
}

func getSafe(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func slugify(s string) string {
	slug := strings.ToLower(s)
	slug = strings.ReplaceAll(slug, " ", "-")
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}

func truncate(maxChars int, s string) string {
	if len(s) <= maxChars {
		return s
	}
	return s[:maxChars] + "..."
}

func formatDate(outputFormat, dateStr string) string {
	parsed, err := time.Parse("January 2, 2006", dateStr)
	if err != nil {
		parsed, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return dateStr
		}
	}
	return parsed.Format(outputFormat)
}

func getFieldValue(fm map[string]interface{}, field string) string {
	if val, ok := fm[field]; ok {
		switch v := val.(type) {
		case string:
			return v
		case int:
			return fmt.Sprintf("%d", v)
		case float64:
			return fmt.Sprintf("%f", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

func getNumericFieldValue(fm map[string]interface{}, field string) (float64, bool) {
	val, ok := fm[field]
	if !ok {
		return 0, false
	}

	switch v := val.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case string:
		n, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err == nil {
			return n, true
		}
	}

	return 0, false
}

func matchesField(post model.Post, field, value string) bool {
	if slice, ok := post.Frontmatter[field].([]interface{}); ok {
		for _, item := range slice {
			if str, isStr := item.(string); isStr && str == value {
				return true
			}
		}
		return false
	}

	return getFieldValue(post.Frontmatter, field) == value
}
