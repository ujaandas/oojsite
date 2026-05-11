---
title: Components
order: 9
template: docs
---

Components are small, reusable HTML fragments in your `components/` directory. You include them in templates and pages using Go's standard templating syntax.

## Basic Component

A component is just HTML:

```html
<!-- components/header.html -->
<header class="site-header">
  <a href="/">Home</a>
  <nav>
    <a href="/about/">About</a>
    <a href="/blog/">Blog</a>
  </nav>
</header>
```

## Including Components

Use Go's `template` action:

```html
{{ template "header.html" . }}
```

Pass `.` to give the component access to post and page data.

## Components with Data

Components receive the same data as their parent template. A component can access `.Frontmatter`, `.Content`, `.Global.Posts`, etc.

**Example: Navigation that highlights the current page**

```html
<!-- components/nav.html -->
<nav>
  {{ range .Global.Posts }}
    {{ if eq .Filepath $.Filepath }}
      <a class="active" href="{{ .Filepath }}">{{ get .Frontmatter "title" }}</a>
    {{ else }}
      <a href="{{ .Filepath }}">{{ get .Frontmatter "title" }}</a>
    {{ end }}
  {{ end }}
</nav>
```

This loops through all posts and marks the current one as active.

## Common Components

**Header:**

```html
<!-- components/header.html -->
<header>
  <h1><a href="/">My Site</a></h1>
  {{ template "nav.html" . }}
</header>
```

**Footer:**

```html
<!-- components/footer.html -->
<footer>
  <p>&copy; 2024 My Name. Built with oojsite.</p>
</footer>
```

**Sidebar:**

```html
<!-- components/sidebar.html -->
<aside>
  <h3>Recent Posts</h3>
  <ul>
    {{ range first 5 (sortByDesc "date" .Global.Posts) }}
      <li><a href="{{ .Filepath }}">{{ get .Frontmatter "title" }}</a></li>
    {{ end }}
  </ul>
</aside>
```

**Tag List:**

```html
<!-- components/tags.html -->
{{ if get .Frontmatter "tags" }}
  <div class="tags">
    {{ range (get .Frontmatter "tags") }}
      <a href="/tags/{{ slugify . }}">{{ . }}</a>
    {{ end }}
  </div>
{{ end }}
```

## Including Multiple Components

Build complex layouts by composing components:

```html
<!DOCTYPE html>
<html>
<head>
  {{ template "meta.html" . }}
</head>
<body>
  {{ template "header.html" . }}
  
  <main>
    {{ .Content }}
  </main>
  
  {{ template "sidebar.html" . }}
  {{ template "footer.html" . }}
</body>
</html>
```

## Conditional Components

Include components conditionally:

```html
{{ if get .Frontmatter "show_sidebar" }}
  {{ template "sidebar.html" . }}
{{ end }}

{{ if get .Frontmatter "show_comments" }}
  {{ template "comments.html" . }}
{{ end }}
```

## Nesting Components

Components can include other components:

```html
<!-- components/sidebar.html -->
<aside>
  {{ template "search.html" . }}
  {{ template "recent-posts.html" . }}
  {{ template "categories.html" . }}
</aside>
```

## Using Template Functions

Components can use all template functions:

```html
<!-- components/post-grid.html -->
{{ range first 6 (sortByDesc "date" (filter "status" "published" .Global.Posts)) }}
  <div class="post-card">
    <h3>{{ get .Frontmatter "title" }}</h3>
    <p>{{ .Snippet }}</p>
  </div>
{{ end }}
```

This shows the 6 most recent published posts.

## Real-World Example

Here's a comment box component:

```html
<!-- components/comments.html -->
{{ if get .Frontmatter "comments" }}
  <section class="comments">
    <h3>Comments</h3>
    <script>
      // Disqus or other comment provider
      var disqus_config = function () {
        this.page.url = "{{ .Filepath }}";
        this.page.identifier = "{{ .SourcePath }}";
      };
    </script>
    <div id="disqus_thread"></div>
  </section>
{{ end }}
```

Then in your template:

```html
<article>
  {{ .Content }}
</article>

{{ template "comments.html" . }}
```

## Best Practices

**Keep components focused** - A component should do one thing well.

**Pass context explicitly** - Always pass `.` to give components full data access.

**Use descriptive names** - `post-meta.html` is clearer than `info.html`.

**Don't repeat markup** - If you're copying HTML between files, make it a component.

**Make them reusable** - Components used in only one place can probably stay inline.
