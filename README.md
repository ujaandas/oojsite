# oojsite

A lightweight, unopinionated static site generator written in Go. It turns Markdown, HTML templates, and reusable components into a static site with no required frontmatter and no custom templating language.

## Quick Start

```sh
nix run github:ujaandas/oojsite -- --allDir docs --dev
```

Then open `http://localhost:8000`.

This runs the docs site directly from the repository. No local setup or state is needed.

## What it handles

- Markdown posts with optional, fully flexible frontmatter
- HTML pages, templates, and reusable components
- Global post collections for indexes, archives, and navigation
- Template helpers for sorting, filtering, grouping, and safe field access
- Optional TailwindCSS output and static file copying

## Layout

```text
posts/       Markdown content
site/        HTML pages
templates/   Post layouts
components/  Shared fragments
static/      CSS, images, scripts
out/         Generated site
```

## More

The full guide lives in [`docs/`](docs/).