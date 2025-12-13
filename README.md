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

- You can insert any post into any page template by using your blog post "tags" frontmatter (i.e; tags: ["posts"] -> {{ .posts }})
- Files you want in your output **MUST BE TRACKED!!!** (ie; `git add site/blog.html`, but not `out/blog.html`) - if you add files while running the watchdog, stop it, `git add` them, and rerun it
- You can use TailwindCSS or just regular CSS rules
- Your public/static content (ie; scripts, CSS, assets, etc...) go in `public/`
- Any templates you want for your blog posts, write them in `templates/`, and use the "template" frontmatter to match it
- Blog posts are recognized as all `*.md` files, and pages are recognized as all `*.html` files

Once you've written whatever content you wanted, just hit `nix run` - it will take care of the rest and give you your built website in `out/`.

> Hint: Use `nix run .#watch-(no-)rb` to track changes across your content for live-reloads. Use no-rebuild mode ONLY if you are changing HTML/CSS exclusively. Changes to files (ie; new posts), or Go source code demands a rebuild.

## 📋 TODO

- Add partials support
- Add Docker support
- Add CI/CD (ie; Github Actions)

## 🐛 Known Bugs

- Sometimes, for no apparent reason, it builds _really_ fast (feature, not a bug, perhaps?)
- If you make and save changes too quickly, your Nix cache might get locked (ie; `error: SQLite database ... is busy`)
- The above will happen, but Nix will simply ignore (?) it
- Sometimes, permissions on `out/` may get garbled up - just delete it and re-run to fix it
