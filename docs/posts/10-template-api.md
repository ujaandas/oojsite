---
title: Template API
order: 10
template: docs
---

This is the complete reference for what's available in your templates and pages.

## Available Data

In any template or page, you have access to a context object:

### For Pages (PageData)

```go
PageData {
  Global: GlobalData
}
```

**`.Global`** - Global data object with all posts

```html
{{ range .Global.Posts }}
  <h2>{{ get .Frontmatter "title" }}</h2>
{{ end }}
```

### For Posts with Templates (TemplateData)

```go
TemplateData {
  Content:    string      // Rendered HTML
  Frontmatter: map[string]interface{}  // YAML fields
  Global:     GlobalData
}
```

**`.Content`** - Converted HTML from Markdown

```html
<article>
  {{ .Content }}
</article>
```

**`.Frontmatter`** - YAML fields from post. Always use `get` to access:

```html
{{ get .Frontmatter "title" }}
{{ get .Frontmatter "author" }}
```

**`.Global.Posts`** - All posts in the site

```html
{{ range .Global.Posts }}
  <p>{{ get .Frontmatter "title" }}</p>
{{ end }}
```

## Post Properties

Each post in `.Global.Posts` has:

```go
Post {
  SourcePath   string                 // Input path (e.g., posts/my-post.md)
  OutputRel    string                 // Output path relative to out/
  Filepath     string                 // Output URL (e.g., /posts/my-post.html)
  Frontmatter  map[string]interface{} // YAML metadata
  Snippet      string                 // First 200 characters (auto-generated)
  Raw          string                 // Original Markdown source
}
```

## Template Functions

All functions are available in templates and pages.

### Data Functions

**`groupBy <field> <posts>`**

Groups posts by a frontmatter field value.

```html
{{ range groupBy "tags" .Global.Posts }}
  <section>
    <h2>{{ .Title }}</h2>
    {{ range . }}
      <p>{{ get .Frontmatter "title" }}</p>
    {{ end }}
  </section>
{{ end }}
```

Returns a slice of groups. Each group has `.Title` (the field value) and `.Items` (the posts).

**`sortBy <field> <posts>`**

Sort posts ascending by a field. Intelligently handles date strings.

```html
{{ range sortBy "date" .Global.Posts }}
  <h3>{{ get .Frontmatter "title" }}</h3>
{{ end }}
```

**`sortByDesc <field> <posts>`**

Sort posts descending. Most recent first for dates.

```html
{{ range sortByDesc "date" .Global.Posts }}
  <h3>{{ get .Frontmatter "title" }}</h3>
{{ end }}
```

**`filter <field> <value> <posts>`**

Filter posts by a field value.

```html
{{ range filter "status" "published" .Global.Posts }}
  <h3>{{ get .Frontmatter "title" }}</h3>
{{ end }}
```

**`first <n> <posts>`**

Get the first N posts.

```html
{{ range first 5 (sortByDesc "date" .Global.Posts) }}
  <h2>{{ get .Frontmatter "title" }}</h2>
{{ end }}
```

**`reverse <posts>`**

Reverse post order.

```html
{{ range reverse (sortBy "date" .Global.Posts) }}
  <p>{{ get .Frontmatter "title" }}</p>
{{ end }}
```

**`unique <field> <posts>`**

Get unique posts by field (removes duplicates).

```html
{{ range unique "author" .Global.Posts }}
  <p>By {{ get .Frontmatter "author" }}</p>
{{ end }}
```

### String Functions

**`get <map> <key>`**

Safely access map values. Returns empty string if missing.

```html
{{ get .Frontmatter "title" }}
{{ get .Frontmatter "optional_field" }}  <!-- Safe if missing -->
```

Always use `get` instead of direct map access.

**`slugify <string>`**

Convert text to URL-safe slug (lowercase, hyphens, no special chars).

```html
<a href="/tags/{{ slugify (get .Frontmatter "tag") }}">
  {{ get .Frontmatter "tag" }}
</a>
```

**`truncate <maxChars> <string>`**

Truncate string to N characters and add `…`

```html
<p>{{ truncate 100 .Snippet }}</p>
```

**`formatDate <format> <dateString>`**

Parse and format dates. Format is Go's `2006-01-02` style.

```html
{{ formatDate "2 Jan 2006" (get .Frontmatter "date") }}
```

Input: `2024-01-15`  
Output: `15 Jan 2024`

## Standard Go Template Functions

All standard Go template functions are available:

**`range`** - Iterate

```html
{{ range .Global.Posts }}
  <p>{{ get .Frontmatter "title" }}</p>
{{ end }}
```

**`if`** - Conditionals

```html
{{ if eq (get .Frontmatter "draft") "true" }}
  <p>This is a draft</p>
{{ else }}
  <p>Published</p>
{{ end }}
```

**`with`** - Conditional with new context

```html
{{ with get .Frontmatter "author" }}
  <p>By {{ . }}</p>
{{ end }}
```

**`len`** - Length

```html
<p>{{ len .Global.Posts }} posts total</p>
```

**`eq`, `ne`, `lt`, `le`, `gt`, `ge`** - Comparisons

```html
{{ if gt (len .Global.Posts) 10 }}
  <p>Many posts!</p>
{{ end }}
```

**`and`, `or`, `not`** - Logical operators

```html
{{ if and (eq .Filepath "/") (gt (len .Global.Posts) 5) }}
  Show something
{{ end }}
```

## Practical Examples

**Blog archive by year:**

```html
{{ range groupBy "year" (sortByDesc "date" .Global.Posts) }}
  <section>
    <h2>{{ .Title }}</h2>
    <ul>
      {{ range .Items }}
        <li>
          <a href="{{ .Filepath }}">{{ get .Frontmatter "title" }}</a>
          <time>{{ get .Frontmatter "date" }}</time>
        </li>
      {{ end }}
    </ul>
  </section>
{{ end }}
```

**Related posts by tag:**

```html
{{ with get .Frontmatter "tags" | first 1 }}
  {{ range filter "tags" . $.Global.Posts }}
    {{ if ne .Filepath $.Filepath }}
      <article>
        <h3>{{ get .Frontmatter "title" }}</h3>
      </article>
    {{ end }}
  {{ end }}
{{ end }}
```

**Recent posts sidebar:**

```html
<aside class="sidebar">
  <h3>Latest</h3>
  {{ range first 5 (sortByDesc "date" .Global.Posts) }}
    <a href="{{ .Filepath }}">{{ truncate 30 (get .Frontmatter "title") }}</a>
  {{ end }}
</aside>
```

**Author page:**

```html
<h1>Posts by {{ get .Frontmatter "author" }}</h1>
{{ range filter "author" (get .Frontmatter "author") .Global.Posts }}
  <article>
    <h2>{{ get .Frontmatter "title" }}</h2>
    <p>{{ .Snippet }}</p>
  </article>
{{ end }}
```

## Piping

Go templates support piping, which chains functions:

```html
<!-- Get the first 10 posts, sorted by date descending -->
{{ range first 10 (sortByDesc "date" .Global.Posts) }}
  <h3>{{ get .Frontmatter "title" }}</h3>
{{ end }}
```

**Complex pipe:**

```html
{{ range first 3 (sortByDesc "date" (filter "published" "true" .Global.Posts)) }}
  <h3>{{ truncate 50 (get .Frontmatter "title") }}</h3>
{{ end }}
```

## Tips

**Always use `get` for safe access** - Trying to access a missing field crashes the template. `get` returns empty string instead.

**Check field exists before using:**

```html
{{ if get .Frontmatter "date" }}
  <time>{{ formatDate "2 Jan 2006" (get .Frontmatter "date") }}</time>
{{ end }}
```

**Use pipes to compose** - Combine functions for powerful queries:

```html
{{ range first 5 (sortByDesc "date" (filter "status" "published" .Global.Posts)) }}
```

**Debug with `{{ . }}`** - Output the current value to see what you're working with.

**Refer to parent scope with `$`** - Inside nested loops, use `$` to reference outer scope:

```html
{{ range .Global.Posts }}
  {{ range .Tags }}
    <a href="/tags/{{ . }}?author={{ $.Author }}">{{ . }}</a>
  {{ end }}
{{ end }}
```
