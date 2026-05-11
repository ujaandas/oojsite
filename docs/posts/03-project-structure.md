---
title: Project Structure
order: 3
template: docs
---

The recommended directory structure is:

```
mysite/
├── posts/          # Markdown content (blog posts, articles)
├── site/           # HTML pages (index, about, archives)
├── templates/      # HTML layouts for posts
├── components/     # Reusable HTML fragments
├── static/         # Images, CSS, scripts
└── out/            # Generated site (created by oojsite)
```

## Directory Breakdown

### `posts/`

Markdown files that become blog posts or articles. Each file can have optional YAML frontmatter.

```markdown
---
title: My Post
date: 2024-01-15
tags: [golang, tutorial]
---

# My Post

Content here...
```

Files can have any structure: flat, nested, or mixed. Subdirectories are preserved in the output.

### `site/`

HTML template files that become pages on your site. These receive the global post collection and can build indexes, archives, or custom layouts.

```html
<!DOCTYPE html>
<html>
  <body>
    {{ range sortByDesc "date" .Global.Posts }}
      <h2>{{ get .Frontmatter "title" }}</h2>
    {{ end }}
  </body>
</html>
```

### `templates/`

HTML layouts that posts can reference in their frontmatter. Use the `template` field to specify which layout a post should use.

```markdown
---
title: My Post
template: article
---
```

Then oojsite will apply `templates/article.html` to wrap your post content.

### `components/`

Small reusable HTML fragments that you include in templates and pages. Include them with standard Go templating:

```html
{{ template "header.html" . }}
{{ template "footer.html" . }}
```

### `static/`

Any static files (images, CSS, JavaScript) that should be copied to the output. If you have `styles.css`, oojsite will optionally run Tailwind CSS on it.

### `out/`

The generated site. This is created by oojsite and should be .gitignored. Each output follows the input structure—posts become `.html` files, pages are rendered, components and templates aren't output directly.

## Flexibility

You don't need all these directories. Some projects:
- Omit templates and write every post with inline styles
- Skip pages and just generate a post feed
- Have only one template for all posts
- Don't use components at all

oojsite adapts to whatever structure you choose.
