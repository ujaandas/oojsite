---
title: Configuration
order: 5
template: docs
---

oojsite uses command-line flags for configuration. There's no config file; everything is passed as arguments.

## Required Paths

These directories must exist:

```bash
oojsite \
  --postDir="posts" \
  --pageDir="site" \
  --templateDir="templates" \
  --componentDir="components" \
  --staticDir="static"
```

If a directory doesn't exist, oojsite will create it.

## Output Settings

**`--outDir`** - Where to write the generated site (default: `out/`)

```bash
oojsite --postDir="posts" --outDir="build"
```

**`--baseURL`** - Base URL for site links (default: `/`)

Use this if your site isn't at the domain root:

```bash
oojsite --baseURL="https://example.com/blog/"
```

This affects how internal links are generated in your templates.

## Development Mode

**`--dev`** - Run a development server on port 8000

```bash
oojsite --postDir="posts" --pageDir="site" --dev
```

The development server watches for changes and rebuilds automatically. All output appears in `out/` (or your specified `--outDir`).

## The allDir Shortcut

Instead of specifying every directory, use `--allDir` to set a root prefix:

```bash
oojsite --allDir docs
```

This is equivalent to:

```bash
oojsite \
  --postDir="docs/posts" \
  --pageDir="docs/site" \
  --templateDir="docs/templates" \
  --componentDir="docs/components" \
  --staticDir="docs/static"
```

If you then override a specific directory, it takes precedence:

```bash
oojsite --allDir docs --postDir="blog/posts"
```

Now posts come from `blog/posts/`, but other directories are still under `docs/`.

## Real-World Examples

**Simple blog:**

```bash
oojsite --allDir myblog --dev
```

**Multi-directory project:**

```bash
oojsite \
  --postDir="content/blog" \
  --pageDir="content/pages" \
  --templateDir="layouts" \
  --componentDir="layouts/components" \
  --staticDir="assets" \
  --outDir="dist"
```

**With custom base URL:**

```bash
oojsite \
  --allDir docs \
  --baseURL="https://mysite.com/docs/" \
  --outDir="public"
```

## Defaults

If you don't specify any directories, oojsite uses:

```
--postDir="posts"
--pageDir="site"
--templateDir="templates"
--componentDir="components"
--staticDir="static"
--outDir="out"
--baseURL="/"
```

So you can start with just:

```bash
oojsite --dev
```

And it will create and use the default structure.
