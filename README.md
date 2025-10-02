# 🐾 oojsite

A tiny static blog generator written in Go.  
Markdown in, HTML out - no fuss (except for me, who had to deal with Nix and Tailwind not playing very nicely with one another).

## 🗂️ Project Structure

- `public/ ` - Static assets (scripts, CSS, UI libs, etc.)
- `templates/` - HTML templates for blog posts and pages
- `content/` - Your actual content:
  - Blog posts: all `*.md` files
  - Pages: all `*.html` files
- `out/` - The generated site output

## 📦 How it Works

I built `oojsite` to be as simple as possible. No DSL or weird templating syntax (ahem, ahem, _Jekyll_).

- You can use TailwindCSS
- Your public/static content (ie; scripts, CSS, UI libs, etc...) go in `public/`
- Any templates you want for your blog posts or pages, write them in `templates/` (more details below)
- Blog posts are recognized as all `*.md` files and pages are recognized as all `*.html` files
- All frontmatter should match and be the same, across both pages and blog posts (keep it simple, stupid)

Once you've written whatever content you wanted, just hit `nix run` - it will take care of the rest and give you your built website in `out/`.

> Hint: Use `nix run .#watch` to track changes across your content for live-reloads.

##
