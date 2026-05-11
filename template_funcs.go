package main

import (
	"fmt"
	"html/template"
	"regexp"
	"sort"
	"strings"
	"time"
)

/*
Template functions provide flexible data manipulation within templates.
All functions are designed to work with user-defined, flexible frontmatter.

Available functions:
  - groupBy: Group items by a frontmatter field
  - sortBy: Sort items by a frontmatter field (ascending)
  - sortByDesc: Sort items by a frontmatter field (descending)
  - filter: Filter items where field equals value
  - first: Get first N items
  - reverse: Reverse item order
  - unique: Get unique values from a field across all items
  - get: Safely access map value (returns empty string if missing)
  - slugify: Convert string to URL-safe slug
  - truncate: Truncate string to N chars with ellipsis
  - formatDate: Parse and format date string (input: "January 2, 2006", output: custom format)
*/

func createTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// Data manipulation
		"groupBy":    groupBy,
		"sortBy":     sortBy,
		"sortByDesc": sortByDesc,
		"filter":     filterPosts,
		"first":      first,
		"reverse":    reversePosts,
		"unique":     unique,

		// Safe accessors
		"get": getSafe,

		// String operations
		"slugify":    slugify,
		"truncate":   truncate,
		"formatDate": formatDate,
	}
}

// groupBy groups posts by a frontmatter field value.
// Returns a map[string][]Post where keys are field values.
// Example: {{ range groupBy .Global.Posts "tags" }}
func groupBy(field string, posts []Post) map[string][]Post {
	result := make(map[string][]Post)
	for _, post := range posts {
		val := getFieldValue(post.Frontmatter, field)
		if val == "" {
			continue
		}

		// Handle both single values and slices
		if slice, ok := post.Frontmatter[field].([]interface{}); ok {
			for _, item := range slice {
				if str, isStr := item.(string); isStr {
					result[str] = append(result[str], post)
				}
			}
		} else {
			result[val] = append(result[val], post)
		}
	}
	return result
}

// sortBy sorts posts by a frontmatter field in ascending order.
// Example: {{ sortBy .Global.Posts "date" }}
func sortBy(field string, posts []Post) []Post {
	sorted := make([]Post, len(posts))
	copy(sorted, posts)

	sort.SliceStable(sorted, func(i, j int) bool {
		valI := getFieldValue(sorted[i].Frontmatter, field)
		valJ := getFieldValue(sorted[j].Frontmatter, field)

		// Try to parse as dates first
		dateI, errI := time.Parse("January 2, 2006", valI)
		dateJ, errJ := time.Parse("January 2, 2006", valJ)
		if errI == nil && errJ == nil {
			return dateI.Before(dateJ)
		}

		// Fall back to string comparison
		return valI < valJ
	})

	return sorted
}

// sortByDesc sorts posts by a frontmatter field in descending order.
// Example: {{ sortByDesc .Global.Posts "date" }}
func sortByDesc(field string, posts []Post) []Post {
	sorted := sortBy(field, posts)
	reversePosts(sorted)
	return sorted
}

// filterPosts filters posts where field equals value.
// Example: {{ filter .Global.Posts "status" "published" }}
func filterPosts(field, value string, posts []Post) []Post {
	var result []Post
	for _, post := range posts {
		if getFieldValue(post.Frontmatter, field) == value {
			result = append(result, post)
		}
	}
	return result
}

// first returns the first N items from a slice.
// Example: {{ first 5 .Global.Posts }}
func first(n int, posts []Post) []Post {
	if n < 0 || n > len(posts) {
		return posts
	}
	return posts[:n]
}

// reversePosts reverses the order of posts.
// Example: {{ reverse .Global.Posts }}
func reversePosts(posts []Post) []Post {
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}
	return posts
}

// unique returns unique values from a frontmatter field across all posts.
// Example: {{ range unique .Global.Posts "tags" }}
func unique(field string, posts []Post) []string {
	seen := make(map[string]bool)
	var result []string

	for _, post := range posts {
		// Handle both single values and slices
		if slice, ok := post.Frontmatter[field].([]interface{}); ok {
			for _, item := range slice {
				if str, isStr := item.(string); isStr && !seen[str] {
					seen[str] = true
					result = append(result, str)
				}
			}
		} else {
			val := getFieldValue(post.Frontmatter, field)
			if val != "" && !seen[val] {
				seen[val] = true
				result = append(result, val)
			}
		}
	}

	return result
}

// getSafe safely retrieves a value from a map, returning empty string if missing.
// Example: {{ get .Frontmatter "title" }}
func getSafe(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// slugify converts a string to a URL-safe slug.
// Example: {{ slugify "My Blog Post" }} -> "my-blog-post"
func slugify(s string) string {
	// Convert to lowercase
	slug := strings.ToLower(s)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove non-alphanumeric characters (except hyphens)
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")
	// Remove consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	// Trim hyphens from start/end
	slug = strings.Trim(slug, "-")
	return slug
}

// truncate truncates a string to maxChars characters with ellipsis.
// Example: {{ truncate 50 "Long text here" }} -> "Long text here..."
func truncate(maxChars int, s string) string {
	if len(s) <= maxChars {
		return s
	}
	return s[:maxChars] + "..."
}

// formatDate parses a date in "January 2, 2006" format and reformats it.
// Example: {{ formatDate "2006-01-02" .Frontmatter.date }}
func formatDate(outputFormat, dateStr string) string {
	parsed, err := time.Parse("January 2, 2006", dateStr)
	if err != nil {
		// Try common formats
		parsed, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return dateStr // Return original if unparseable
		}
	}
	return parsed.Format(outputFormat)
}

// Helper: safely get a string value from frontmatter map
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
