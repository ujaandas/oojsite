---
title: Pages
order: 6
template: docs
---

Pages are HTML template files in your `site/` directory. They have access to all posts and build the structural pages of your site: home page, archive, about page, etc.

## Basic Page

A page is a simple HTML template:

```html
<!DOCTYPE html>
<html>
<head>
  <title>My Site</title>
</head>
<body>
  <h1>Welcome</h1>
  <p>This is my site.</p>
</body>
</html>
```

When you run oojsite, this file becomes `out/index.html` (renamed from `index.html`).

## Accessing Posts

Pages receive a `PageData` object with access to all posts:

```html
{{ range .Global.Posts }}
  <h2>{{ get .Frontmatter "title" }}</h2>
{{ end }}
```

This is how you build post archives, galleries, or any dynamic index.

## Common Patterns

**Recent Posts Archive:**

```html
{{ range first 10 (sortByDesc "date" .Global.Posts) }}
  <article>
    <h2>{{ get .Frontmatter "title" }}</h2>
    <p>{{ get .Frontmatter "date" }}</p>
    <a href="{{ .Filepath }}">Read more</a>
  </article>
{{ end }}
```

**Posts by Tag:**

```html
{{ range groupBy "tags" .Global.Posts }}
  <section>
    <h2>Tag: {{ .Title }}</h2>
    {{ range . }}
      <p>{{ get .Frontmatter "title" }}</p>
    {{ end }}
  </section>
{{ end }}
```

**Filter by Status:**

```html
{{ range filter "status" "published" .Global.Posts }}
  <h3>{{ get .Frontmatter "title" }}</h3>
{{ end }}
```

## Including Components

Use Go's standard template syntax to include components:

```html
{{ template "header.html" . }}

<main>
  <!-- Your content -->
</main>

{{ template "footer.html" . }}
```

## Using Page-Specific Data

You can add frontmatter to your HTML file using HTML comments at the top (though this is less common):

```html
<!--
title: My Page
description: A description for meta tags
-->
<!DOCTYPE html>
...
```

However, pages typically don't use frontmatter—they're mostly static structure with dynamic post lists.

## Multiple Pages

Create multiple HTML files in `site/`:

```
site/
├── index.html        → out/index.html
├── archive.html      → out/archive.html
├── about.html        → out/about.html
└── projects/
    └── index.html    → out/projects/index.html
```

Directory structure is preserved in the output.

## Advanced Example

```html
<!DOCTYPE html>
<html>
<head>
  <title>Archive</title>
</head>
<body>
  <h1>All Posts</h1>
  
  {{ range groupBy "year" (sortByDesc "date" .Global.Posts) }}
    <section>
      <h2>{{ .Title }}</h2>
      <ul>
        {{ range . }}
          <li>
            <a href="{{ .Filepath }}">{{ get .Frontmatter "title" }}</a>
            <span>{{ get .Frontmatter "date" }}</span>
          </li>
        {{ end }}
      </ul>
    </section>
  {{ end }}
</body>
</html>
```

This groups posts by year and creates an organized archive page.
