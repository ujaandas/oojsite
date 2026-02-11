---
title: "Streamlining my Resume with Nix and LaTex"
template: "post"
tags: ["projects"]
random: "pepee"
---

One of the biggest problems I had with my resume previously was fiddling between different layouts for different positions (i.e; full-stack oriented resume, backend oriented, research oriented, etc...). So, to solve this, I spent the past days transforming my Nix-backed LaTeX resume builder into a fully modular resume generator. Now I can write each content 'block' exactly once (i.e; a single source of truth) and assemble it into any layout I want with Nix.

If you just want the code, check out the [repo](https://github.com/ujaandas/nix-resume).

### Defining Variants and Templates

I organize variants as Nix attribute sets. Each variant specifies a template file and a map of section names to content files:

```nix
variants = {
  a = {
    template = ./src/templates/a.tex;
    sections = {
      "Work Experience" = [
        ./src/content/work/work1.tex
        ./src/content/work/work2.tex
      ];
      "Education" = [
        ./src/content/edu/education.tex
      ];
    };
  };

  b = {
    template = ./src/templates/b.tex;
    sections = {
      "Projects" = [
        ./src/content/projects/proj1.tex
        ./src/content/projects/proj2.tex
      ];
      "Skills" = [
        ./src/content/skills/skills.tex
      ];
    };
  };
};
```

By mapping each variant to its own template and section list, I avoid duplicating content across multiple `.tex` files.

### Parsing and Building

The core of this system is a small Nix builder in `build.nix`. It reads the template, finds section markers like `\section\*{Work Experience}`, and injects the corresponding `\input{}` lines:

```nix
mkDoc =
    templateStr: sections:
    let
      lines = builtins.split "\n" templateStr;
      sectionNames = builtins.attrNames sections;

      inject =
        line:
        let
          matched = builtins.filter (t:
            hasSubstr
              ("\\\\section\\*\\{" + t + "\\}")
              line)
            sectionNames;
        in
        if matched != [ ] then
          line
          + "\n"
          + builtins.concatStringsSep "\n" (
            map
              (f: "\\input{" + toString f + "}")
              sections.${builtins.head matched}
          )
        else
          line;

      allLines = builtins.map inject lines;
    in
    builtins.concatStringsSep
      "\n"
      (builtins.filter
        builtins.isString
        allLines);
```

This function produces a complete `main.tex` that `latexmk` can compile without manual edits.

And that's pretty much it. Now I can set up some CI/CD to automatically deploy different versions of my resume to whereever I so please. Nice!
