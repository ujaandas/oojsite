# 🐾 oojsite

A lightweight, unopinionated static site generator written in Go. Transform Markdown and HTML templates into a complete static site with zero constraints on content structure.

## Quick Start

```sh
# Try the example (yes, I know, crappy example, will replace soon!)
nix run . -- \
  --pageDir="example/site" \
  --postDir="example/posts" \
  --staticDir="example/static" \
  --templateDir="example/templates" \
  --componentDir="example/components" \
  --dev

# Then visit http://localhost:8000
```

## Features

**Completely Flexible**
- No mandatory frontmatter fields - define whatever you need
- Works for blogs, portfolios, documentation, catalogs, anything
- Pure Go templating, no custom DSL

**Powerful, But Simple Template Functions**
- Group, sort, filter, and manipulate content dynamically
- Safe field access with automatic fallbacks
- String operations (slugify, truncate, date formatting)

**Simple Architecture**
- Markdown -> HTML conversion with YAML frontmatter
- Reusable components and layouts
- Global post collection available in all templates
- TailwindCSS support

All paths are configurable via CLI flags.

## Getting Started

### 1. Write Some Markdown

**Markdown posts** (`posts/my-post.md`):
```yaml
---
title: My First Post
date: January 15, 2024
tags: [ "golang" "tutorial" ]
custom_field: anything you want
---

# Post content in Markdown

Your content here...
```
### 2. Write Some HTML

**HTML pages** (`site/index.html`):
```html
<!DOCTYPE html>
<html>
  <head><title>My Site</title></head>
  <body>
    {{ range sortByDesc "date" .Global.Posts }}
      <h2>{{ get .Frontmatter "title" }}</h2>
    {{ end }}
  </body>
</html>
```

### 3. Generate!

```bash
oojsite --postDir="posts" --pageDir="site" --templateDir="templates" --componentDir="components" --dev
```

## How It Works

All templates have access to:
- **`.Frontmatter`** – User-defined YAML fields
- **`.Content`** – Converted HTML (posts only)
- **`.Global.Posts`** – All posts processed so far
- **Template functions** – groupBy, sortBy, filter, get, slugify, etc.

Essentially, `oojsite` breaks your content into 4 content types:
1. Posts, which are any Markdown files in `postDir`
2. Pages, which are all the HTML files in `pageDir`
3. Templates, which are all the HTML files in `templateDir`
4. Components, which are all the HTML files in `componentDir`

The implementation is intentionally simple:
- Posts may have frontmatter, but they do not need it
- If a post specifies `template: my-layout`, `oojsite` uses `templates/my-layout.html`
- If no template is specified, the markdown is converted straight to HTML and written out directly
- Components are reusable HTML fragments that you can call from pages or templates with `{{ template "header.html" . }}`
- Pages are plain HTML templates that get the global post collection and can build indexes, archives, or anything else

For a new user, the shortest path is:
1. Put Markdown files in `posts/`
2. Put page templates in `site/`
3. Put reusable layouts in `templates/` and shared fragments in `components/`
4. Run `oojsite --dev`

That is enough to get started. The only real convention is how you choose to structure your own frontmatter and how you use the template helpers.



## Template API

### Accessing Frontmatter

Always use `get` to safely access fields (returns empty string if missing):

```html
{{ get .Frontmatter "title" }}
{{ get .Frontmatter "optional_field" }}  <!-- Safe if field doesn't exist -->
```

### Data Functions

**`groupBy <field> <posts>`** – Group by any frontmatter field
```html
{{ range groupBy "tags" .Global.Posts }}
  {{ .Title }}: {{ . | len }} items
{{ end }}
```

**`sortBy <field> <posts>`** – Sort ascending (handles dates intelligently)
```html
{{ range sortBy "date" .Global.Posts }}
  {{ get .Frontmatter "title" }}
{{ end }}
```

**`sortByDesc <field> <posts>`** – Sort descending
```html
{{ range sortByDesc "date" .Global.Posts }}
  {{ get .Frontmatter "title" }}
{{ end }}
```

**`filter <field> <value> <posts>`** – Filter by field value
```html
{{ range filter "status" "published" .Global.Posts }}
  {{ get .Frontmatter "title" }}
{{ end }}
```

**`first <n> <posts>`** – Get first N items
```html
{{ range first 5 (sortByDesc "date" .Global.Posts) }}
  {{ get .Frontmatter "title" }}
{{ end }}
```

**`reverse <posts>`** – Reverse order
```html
{{ range reverse .Global.Posts }}
  {{ get .Frontmatter "title" }}
{{ end }}
```

**`unique <field> <posts>`** – Get unique values from a field
```html
{{ range unique "tags" .Global.Posts }}
  <a href="/tags/{{ slugify . }}">{{ . }}</a>
{{ end }}
```

### String Functions

**`slugify <string>`** – Convert to URL-safe slug
```html
{{ slugify "My Blog Post" }}  <!-- "my-blog-post" -->
```

**`truncate <maxChars> <string>`** – Truncate with ellipsis
```html
{{ truncate 50 .Snippet }}
```

**`formatDate <format> <dateString>`** – Parse and reformat dates
```html
{{ formatDate "2006-01-02" (get .Frontmatter "date") }}
```

**`get <data> <key>`** – Safe map access
```html
{{ get .Frontmatter "title" }}
```