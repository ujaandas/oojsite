---
title: Why oojsite?
order: 1
template: docs
---

Most static site generators impose structure: required frontmatter fields, rigid content hierarchies, or custom template syntax. oojsite exists to avoid all that.

## The Problem

You might need:
- A blog with posts, but also a portfolio with projects
- Tags, categories, and custom metadata
- Some posts with custom layouts, others with plain HTML
- Flexibility to change your content model as your needs evolve

Traditional generators force you into their opinion. You end up with unused fields, workarounds, and templates that don't quite fit.

## The Solution

oojsite is **intentionally unopinionated**. It provides:

- **Zero mandatory frontmatter** - Define only the fields you need. A post can be just Markdown with no YAML at all.
- **Flexible content types** - Posts, pages, templates, components. Combine them however you want.
- **Simple but powerful** - No custom DSL. Just Go templates and a straightforward function API.
- **Graceful degradation** - Works without Tailwind CSS. Works without templates. Works without even frontmatter.

## Design Philosophy

1. **Content is data** - Your Markdown becomes HTML with accessible frontmatter fields.
2. **Templates are simple** - Use standard Go `html/template` syntax. No learning curve.
3. **Composition over convention** - Build your site's structure through templates and components.
4. **Convention-free** - No magic filenames or directory structures (except what you configure).

## Who Should Use oojsite?

- You want complete control over your content model
- You are comfortable with Go templates
- You need flexibility without boilerplate
- You build sites with Markdown or HTML
- You value simplicity over features
