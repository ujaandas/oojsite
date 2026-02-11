---
title: "Building a Static Site Generator in Go"
template: "post"
tags: ["projects"]
---

I feel like my portfolio sites keep getting dumber and dumber, and maybe that's a good thing. It started with Next.js, which, for obvious reasons (for a static site),
was a really bad idea and totally overkill. Then, I moved to Jekyll. Honestly, I think Jekyll is just fine, and I just like being a bit of a fusspot. I didn't like the
templating syntax , [Liquid](https://jekyllrb.com/docs/liquid/), I didn't like the myriad of random config files with what felt like very arbitrary options, and I didn't
like writing CSS. I felt like it was a simple-enough task for me to do well-enough myself.

### Building Oojsite

I was, for the most part, correct - building `oojsite` was fun. I like Go, I like this middle-ground of backend and web, and I like building tooling I actually use.

So, what did I want? I wanted something _ergonomic_, where all I'd have to do is write my post and hit send. If you're wondering, yes, that is exactly how my previous
site with Jekyll was. Still, I have full-control over the build process here and can reason about the templating, styling and structure of `oojsite` a lot more clearly.

### How's it work?

`oojsite` uses Go's built-in [html/template](https://pkg.go.dev/html/template), meaning if you've used Go templates before, this is literally nothing new. Every `* html` is basically treated like a template - `oojsite` globs these files, and builds a [template](https://pkg.go.dev/text/template) from which we can make the appropriate lookups:

```go
// Load both page templates (ie; for posts) and actual pages (ie; index.html).
func loadPages(tmplDir, siteDir string) (*template.Template, error) {
	tmpls := template.New("")

	// Load post templates
	tmpls, err := tmpls.ParseGlob(fmt.Sprintf("%s/*.html", tmplDir))
	if err != nil {
		return nil, err
	}

	// Load actual page templates
	tmpls, err = tmpls.ParseGlob(fmt.Sprintf("%s/*.html", siteDir))
	if err != nil {
		return nil, err
	}

	return tmpls, nil
}
```

Once we've loaded the templates, we can transpile our Markdown posts into HTML using [Goldmark](https://github.com/yuin/goldmark). Like the rest of the app, I opted for the simplest possible pipeline: read file, split frontmatter, run Markdown through a renderer, hand the resulting HTML off to the template engine, and write the resultant page out:

```go
func processPost(path, postDir string, tmpls *template.Template) error {
	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Extract frontmatter
	post, err := extractFrontmatter(path, content)
	if err != nil {
		return err
	}

	// Convert markdown to HTML
	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert(post.Raw, &buf); err != nil {
		return err
	}

	// Apply template
	tmpl := tmpls.Lookup(fmt.Sprintf("%s.html", post.Frontmatter.Template))
	if tmpl == nil {
		return err
	}

	// Write outputted file
	return writePostFile(path, postDir, tmpl, post, buf)
}
```

Also to note is that the frontmatter is _very_ important:

```go
type Frontmatter struct {
	Title    string   `yaml:"title"`
	Template string   `yaml:"template"`
	Tags     []string `yaml:"tags,omitempty"`
}
```

Most of these make sense, like the title (obviously), the template (matched by name, determines the outline of the post page), but what's interesting are the tags. Tags are unique in that you can call that tag from any other HTML page to iterate over all pages with that tag. This lets me split my posts by category, and maybe later add proper search functionality.

After all the content is generated, we build Tailwind. This just shells out to the Tailwind CLI with the right input/output paths. I didn’t want to embed a CSS pipeline into the tool - that’s how you end up reinventing Webpack by accident - so calling the CLI is perfectly fine (more on that and dependencies later!):

```go
// Use `os/exec` to build TailwindCSS.
func buildTailwind(outDir string) error {
	// TODO: Make these options/user-changeable later
	in := filepath.Join("static", "styles.css")
	out := filepath.Join(outDir, "static", "styles.css")

	// Assumes the user has `tailwindcss` available
	cmd := exec.Command("tailwindcss", "--input", in, "--output", out, "--minify", "--content", fmt.Sprintf("%s/**/*.html", outDir))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
```

Finally, we copy static files. Everything in `static/` is mirrored into the output directory. No processing, bundling, or minification (other than Tailwind, teehee). If you put a file in `static/`, it ends up in `out/static/`.

> I was kind of surprised by how tedious this was. 67 (haha) lines to move a folder? I might as well have just called `cp`.

### Getting Nixy

The last piece of the puzzle is Nix. I wanted the build to be reproducible, portable, and not depend on whatever random version of Go or Tailwind I happened to have installed. The flake for this site used to be a lot more complicated as I would build Tailwind in Nix, but have since delegated that to Go:

```nix
oojsite = pkgs.buildGoModule {
    pname = "oojsite";
    version = "0.1.0";
    src = ./.;
    vendorHash = "sha256-gM37SLXNi4uY3uetmagNarbUvaFapQciajrguWVSd34=";
    buildInputs = with pkgs; [
    makeWrapper
    tailwindcss
    ];

    postInstall = ''
    wrapProgram $out/bin/oojsite \
        --prefix PATH : "${pkgs.tailwindcss}/bin"
    '';

};
```

All this does is:

- Build `oojsite` using the default Go package builder for the current system
- Wrap the `tailwindcli` around it (janky, I know)
- Expose the built site as `packages.oojsite`
- Make it the default package so `nix run` "just works" (I also have some other funky options using `watchdog`)

So now, to build the entire website - templates, posts, Tailwind, static assets - I literally just run:

```nix
nix run
```

and the whole thing pops out in `out/`. Oh, you'd like to change the output folder? Well, on that note...

### Future Development

I have a lot to do here! Firstly, I do have a working option configuration system, I just haven't really tested it (or anything else, for that matter). Also, not everyone likes/uses Nix, so Docker would be nice. Oh, and if you've read the source code, you might notice that [currently](https://github.com/ujaandas/ujaandas.github.io/commit/bff08295d3733ef1f5aa11cc6dfd9507a6300f6c),
I repeat a bunch of code (what happened to DRY?), so partials would also be neat. There's a few other things too, better ergonomics, documentation, unit-testing like I mentioned before, etc...

But, even still, I think `oojsite` is a handy-dandy little tool that I'm happy to keep iterating on (for now).
