# 🐾 oojsite

A tiny and ergonomic static blog generator written in Go.  
Markdown in, HTML out - no fuss.

## Example
Check out my own [blog](https://github.com/ujaandas/ujaandas.github.io/) for an example. You can also use the `flake.nix` there as a starter template.

## Project Structure

- `static/ ` - Static assets (scripts, CSS, UI libs, etc.)
- `templates/` - HTML templates for blog posts
- `site/` - HTML templates for actual pages
- `posts/` - Blog posts and contents in Markdown
- `out/` - The generated site output

> All of these folder names can be changed! Check `options.go`.

## How it Works

I built `oojsite` to be as simple as possible. No DSL or weird templating syntax (ahem, ahem, _Jekyll_), it's just Go's text/template.

- You can insert any post into any page template by using your blog post "tags" frontmatter (i.e; tags: ["posts"] -> {{ .posts }})
- You can use TailwindCSS or just regular CSS rules (for instance, in `styles.css`)
- Your public/static content (ie; scripts, CSS, assets, etc...) goes in `static/`
- Currently, you need a dummy `styles.css`, if you don't want to use Tailwind, just omit the headers
- Any templates you want for your blog posts, write them in `templates/`, and use the "template" frontmatter to match it
- Blog posts are recognized as all `*.md` files, and pages are recognized as all `*.html` files

Once you've written whatever content you wanted, just hit `nix run` - it will take care of the rest and give you your built website in `out/`.

## TODO

- Better ergonomics
- Global variables and template state
- Add partials support
- Add Docker support

## 🐛 Known Bugs

- So far so good?
