---
title: Develop and Build
order: 4
template: docs
---

## Development Mode

Run oojsite in development mode to build your site and start a local server:

```bash
oojsite ... --dev
```

Then visit `http://localhost:8000` to see your site.

Currently, oojsite does _not_  watch for changes pr rebuild automatically. You have to restart to see changes (for now).

## Building for Production

To generate the output without starting a server:

```bash
oojsite ... --outDir="out"
```

The generated site will be in `out/`. You can then deploy this to any static hosting service.

## Using allDir

If your project is organized with a single root directory:

```bash
oojsite --allDir myproject --dev
```

This automatically finds:
- `myproject/posts`
- `myproject/site`
- `myproject/templates`
- `myproject/components`
- `myproject/static`

## Full CLI Options

```
--postDir string
    Directory containing Markdown posts (default "posts")

--pageDir string
    Directory containing HTML page templates (default "site")

--templateDir string
    Directory containing post templates (default "templates")

--componentDir string
    Directory containing reusable components (default "components")

--staticDir string
    Directory containing static files (default "static")

--outDir string
    Output directory for generated site (default "out")

--baseURL string
    Base URL for site links (default "/")

--allDir string
    Convenience prefix for all directories

--dev
    Run development server on :8000
```

## Building with Nix

If you're using Nix, the flake provides a development shell:

```bash
nix flake
oojsite --allDir docs --dev
```

The flake also provides a build output for creating reproducible builds.

## Testing Locally

After building, you can serve the output with any HTTP server:

```bash
cd out
python3 -m http.server 8000
```

Or use `npx`:

```bash
npx serve out
```

Then visit `http://localhost:8000`.

## Deploying

Common deployment targets:

**GitHub Pages:**
- Push your `out/` directory to a `gh-pages` branch
- GitHub automatically serves it at `yourname.github.io`

**Netlify:**
- Connect your repo
- Set build command to your oojsite command
- Set publish directory to `out/`

**Any static host:**
- Copy the `out/` directory to your server
- Any web server (nginx, Apache, S3, Cloudflare, etc.) works
