---
title: "Nix and LaTex"
template: "post"
tags: ["posts"]
random: "pepee"
---

Lately, I've been wondering why I've had to put up with Overleaf's awful UX when I could just as well edit my personal documents on my own code editor and track relevant changes with Git. Unfortunately, the first thought that came to my mind was how I could build said documents with Nix.

### The TeX Ecosystem

Truth be told, I don't know a whole lot about the TeX ecosystem - every research paper or report I've done, I've used Overleaf and offloaded my thinking and ecosystem handling to them. The closest I've come to fiddling with TeX (pdfTeX?) on my local machine was building a pipeline to auto-convert Jupyter Notebooks into PDFs for my machine learning class at Northwestern.

This [website](https://www.tug.org/levels.html) does a _really_ good job at summarizing all things TeX, so here's a summary of the summary:

- TeX comes in various distributions (i.e; complete bundlings of all things TeX) and in my case, I will be using TeX Live
- TeX distributions contain various engines - these target specific outputs, such as pdfTeX for pdf output, or XeTeX for unicode
- TeX comes in a variety of flavours in the form of formats, the most popular being LaTeX
- TeX also comes with several optional packages suited to specific, niche (?) problems

### Building a Flake

Next, building the flake.

TeX distributions, especially Live, can come with thousands of packages, most of which I probably won't need. Luckily, Nix lets us create [subsets](https://nixos.wiki/wiki/TexLive) to choose what we want and don't want.

Initially I tried building off of the small set, combining whatever else I wanted on top of that using the `.combine` directive.

```nix
tex =
  with pkgs;
  texlive.combine {
    inherit (texlive) scheme-small latexmk;
  };
```

But found that my workflows needed packages which needed packages which needed more packages, so I just ended up using scheme-basic.

### First Flake

As it stood, my initial `flake.nix` looked something like this:

```nix
{
  description = "Example development environment flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        tex =
          with pkgs;
          texlive.combine {
            inherit (texlive) scheme-basic latexmk;
          };
      in
      {
        defaultPackage = pkgs.stdenv.mkDerivation {
          name = "latex-document";
          src = ./.;
          nativeBuildInputs = [
            tex
          ];

          buildPhase = ''
            latexmk -pdf main.tex
          '';

          installPhase = ''
            mkdir -p $out
            cp main.pdf $out/
          '';
        };
      }
    );
}
```

Next up, I wanted to create a way to save each section of my document (e.g; each work experience, each project, each publication, etc...) as its own chunk, and piece/order them using Nix.
