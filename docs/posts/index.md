---
title: oojsite
slug: home
order: 0
---

Welcome to **oojsite**, a minimal static site generator that turns Markdown and HTML templates into complete static sites. This documentation site is built with oojsite itself, demonstrating its capabilities.

## What is oojsite?

oojsite is a command-line tool that processes:
- **Markdown files** into blog posts and articles
- **HTML templates** into site pages and layouts
- **Static assets** (CSS, images, scripts)
- **Reusable components** for common UI patterns

The philosophy is simple: **no mandatory structure, no conventions, complete flexibility**.

## Core Concepts

- **Posts** are Markdown files that live in your posts directory. They may have optional YAML frontmatter with any fields you define. If you specify a template, oojsite applies that layout to your post. If not, it converts the Markdown directly to HTML.

- **Pages** are HTML templates that have access to all posts and can build indexes, archives, galleries, or anything else. Think of them as the structural pages of your site.

- **Templates** are reusable HTML layouts that format individual posts. A post's frontmatter can reference which template to use.

- **Components** are small HTML fragments you include in templates and pages via standard Go templating (`{{ template "name.html" . }}`).