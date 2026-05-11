---
title: Templates
order: 8
template: docs
---

Templates are HTML files in your `templates/` directory. When a post specifies a template in its frontmatter, oojsite applies that layout to wrap the post content.

## Basic Template

A template is a standard Go HTML template with access to post data:

```html
<!DOCTYPE html>
<html>
<head>
  <title>{{ get .Frontmatter "title" }}</title>
</head>
<body>
  <header>
    <h1>{{ get .Frontmatter "title" }}</h1>
  </header>
  <article>
    {{ .Content }}
  </article>
  <footer>
    <p>Published {{ get .Frontmatter "date" }}</p>
  </footer>
</body>
</html>
```

Save this as `templates/article.html`, then use it in a post:

```markdown
---
title: My Post
template: article
---

# Content here
```

oojsite will convert the Markdown, then pass `.Content` to the template, which wraps it in your layout.

## Accessing Post Data

Templates have full access to the post object:

- **`.Frontmatter`** - YAML fields (always use `get` for safety)
- **`.Content`** - Rendered HTML from Markdown
- **`.Snippet`** - First 200 characters
- **`.Filepath`** - Output file path
- **`.Global.Posts`** - All posts (useful for navigation)

## Including Components

Use Go's `template` action to include components:

```html
{{ template "header.html" . }}

<article>
  {{ .Content }}
</article>

{{ template "footer.html" . }}
```

Pass `.` to give components access to all post data.

## Conditional Logic

You can use Go template conditionals:

```html
{{ if eq (get .Frontmatter "draft") "true" }}
  <div class="warning">This is a draft.</div>
{{ end }}

{{ if get .Frontmatter "updated" }}
  <p>Updated {{ get .Frontmatter "updated" }}</p>
{{ end }}
```

## Styling with Classes

Apply CSS classes based on frontmatter:

```html
<article class="{{ get .Frontmatter "style" }}">
  {{ .Content }}
</article>
```

Then in your post:

```markdown
---
title: Dark Theme Post
style: article dark-theme
---
```

## Navigation Between Posts

Since `.Global.Posts` is available, you can add previous/next links:

```html
<!-- Simplified example -->
{{ range $i, $post := .Global.Posts }}
  {{ if eq $post.Filepath $.Filepath }}
    {{ if gt $i 0 }}
      {{ with index .Global.Posts (sub $i 1) }}
        <a href="{{ .Filepath }}">← Previous</a>
      {{ end }}
    {{ end }}
    
    {{ if lt $i (sub (len .Global.Posts) 1) }}
      {{ with index .Global.Posts (add $i 1) }}
        <a href="{{ .Filepath }}">Next →</a>
      {{ end }}
    {{ end }}
  {{ end }}
{{ end }}
```

## Multiple Templates

Create different templates for different post types:

```
templates/
├── article.html      # Standard blog post
├── photo.html        # Photo essay
└── snippet.html      # Quick tip
```

Each post chooses which template suits it:

```markdown
---
title: A Beautiful Sunset
template: photo
---
```

## No Template Required

If a post doesn't specify a template, oojsite skips the template layer entirely. It converts Markdown to HTML and writes it directly. This is useful for simple pages that don't need any wrapper.

## Template Functions

Inside templates, you can use all template functions like `sortBy`, `filter`, `groupBy`, etc. See the Template API section for the complete reference.

## Real Example

Here's a practical blog post template:

```html
<!DOCTYPE html>
<html>
<head>
  <title>{{ get .Frontmatter "title" }} | Blog</title>
  <meta name="description" content="{{ get .Frontmatter "summary" }}">
</head>
<body>
  {{ template "navigation.html" . }}
  
  <main>
    <article class="post">
      <header>
        <h1>{{ get .Frontmatter "title" }}</h1>
        <div class="meta">
          <time>{{ get .Frontmatter "date" }}</time>
          {{ if get .Frontmatter "author" }}
            by <span class="author">{{ get .Frontmatter "author" }}</span>
          {{ end }}
        </div>
      </header>
      
      {{ .Content }}
      
      {{ if get .Frontmatter "tags" }}
        <footer class="tags">
          Tags: {{ get .Frontmatter "tags" }}
        </footer>
      {{ end }}
    </article>
    
    {{ template "related-posts.html" . }}
  </main>
  
  {{ template "footer.html" . }}
</body>
</html>
```

This template uses components for navigation, related posts, and footer, keeping each piece modular and reusable.
