---
title: Posts
order: 7
template: docs
---

Posts are Markdown files in your `posts/` directory. Each post becomes an HTML file in the output.

## Basic Post

The simplest post is just Markdown:

```markdown
# My First Post

This is my content. No frontmatter needed.
```

oojsite converts it to HTML and writes it to the output.

## With Frontmatter

Add optional YAML frontmatter to define metadata:

```markdown
---
title: My First Post
date: 2024-01-15
tags: [golang, tutorial]
author: Jane Doe
custom_field: anything you want
---

# My First Post

Content here...
```

**Important:** All frontmatter fields are optional. Only define what you need. There are no required fields.

## Common Fields

While frontmatter is flexible, here are commonly used fields:

- `title` - Post headline
- `date` - Publication date (any format)
- `tags` - List of tags/categories
- `summary` - Short description for previews
- `template` - HTML layout file to use
- `draft` - Mark as draft (use `filter` to exclude)
- `author` - Post author

But you can use any fields you want.

## Specifying a Template

If you want to apply a custom layout to your post, use the `template` field:

```markdown
---
title: My Article
template: article
---

# Content
```

oojsite will apply `templates/article.html` to wrap your post content.

If you don't specify a template, your Markdown is converted directly to HTML with no wrapper.

## No Frontmatter

A post doesn't need frontmatter at all:

```markdown
# Just Markdown

This entire file is treated as body content. There's no metadata, but that's fine if you don't need it.
```

## Post Output

Posts are output as `.html` files, preserving directory structure:

```
posts/
├── my-first-post.md          → out/my-first-post.html
├── 2024/
│   ├── january.md            → out/2024/january.html
│   └── february.md           → out/2024/february.html
```

The file path is available in templates as `.Filepath`.

## Accessing Post Data

In your template, posts have:

- **`.Frontmatter`** - Map of YAML fields. Always use `get` to access: `{{ get .Frontmatter "title" }}`
- **`.Content`** - Converted HTML from Markdown
- **`.Snippet`** - First 200 characters of body (auto-generated)
- **`.Raw`** - Original Markdown source
- **`.Filepath`** - Path to output file (e.g., `/posts/my-post.html`)
- **`.SourcePath`** - Path to input file
- **`.Global.Posts`** - All posts processed so far (available in templates)

## Example Template

```html
<!DOCTYPE html>
<html>
<head>
  <title>{{ get .Frontmatter "title" }}</title>
</head>
<body>
  <header>
    <h1>{{ get .Frontmatter "title" }}</h1>
    <p>By {{ get .Frontmatter "author" }} on {{ get .Frontmatter "date" }}</p>
    {{ if eq (get .Frontmatter "draft") "true" }}
      <span class="draft">Draft</span>
    {{ end }}
  </header>
  <article>
    {{ .Content }}
  </article>
</body>
</html>
```

## Tips

**Use consistent date formats** - If you plan to sort by date in pages, use a consistent format like `YYYY-MM-DD`.

**Keep frontmatter lean** - Only include fields you actually use. Extra frontmatter just adds noise.

**Leverage the snippet** - The auto-generated 200-character snippet is useful for previews. Don't store redundant summaries unless you need custom text.

**Use filter for drafts** - Mark posts as `draft: true` and filter them out in pages:

```html
{{ range filter "status" "published" .Global.Posts }}
  <!-- only published posts -->
{{ end }}
```
