---
title: "Why Nix?"
template: "post"
tags: ["posts"]
random: "pepee"
---

I reckon if you asked any Mac user to list how many programs they had installed via Brew, they'd either give you a blank stare and wonder why you're asking a stranger about beer-making, or give you an answer in the ballpark of a couple hundred.

### Reproducibility

Far too often did I find myself installing a bunch of random tools and programs to test out random stuff, and then promptly forget to ever remove them. And those 5 different versions of Python you installed (and inevitably have to wrangle via Conda)? Five different Java SDKs? Node version manager? The list goes on...

We've all been there. Every developer knows the woes of dependency hell, and arguably the biggest cause is a lack of reproducibility.

That's where Nix comes in. Be it on a system-level, environment-level, or even package-level, Nix is a tool that was developed for just that, and one that I've come to really appreciate in the past few months of using it.

From my understanding, Nix is both a package manager and a language, where Nix the language is used to describe "inputs" and "outputs" to Nix the package manager. Through fully reproducible, cross-platform, binary-cached builds from the largest package repository in the world, Nix allows me to create the fully reproducible system I've always wanted.

### Nix is Pretty Straightforward

1. **Derivations**: Nix packages are defined using a functional language. Each package definition, called a derivation, describes how to build the package, including its dependencies, source code, and build instructions.

2. **Hashes**: Each derivation produces a unique hash based on its inputs (dependencies, configuration options, source code, etc.). This hash uniquely identifies the package's build outputs.

3. **Reproducibility**: Since the hash is generated from all the inputs, if you have the same inputs, you'll get the same hash and build output, regardless of the environment. This ensures that builds are reproducible across different machines.

4. **Isolation**: Nix builds packages in isolation using unique build environments. This prevents interference from other packages or system settings (if in a pure environment), ensuring that the build process is clean and reproducible.

5. **Binary Cache**: Nix also has a global binary cache where pre-built packages are stored. If a package's hash matches one in the cache, Nix can download and use the pre-built binary instead of building it from scratch.

6. **Pinning**: By pinning dependencies to specific hashes, Nix ensures that your projects always use the exact same versions of dependencies, avoiding issues caused by updates or changes in dependencies. This is particularly useful for long-term projects or collaborations where consistency is crucial.

### Trying it out

This is probably Nix's biggest downfall, and I'm sure most of the community would agree. Unlike Arch or somesuch, which is usually as simple as booting in and selecting your settings, Nix seems to lack a proper, easy-to-follow, opinionated guide to set it up and learn it from the get go. As a result, getting Nix up and running can be a little troublesome, but it's still insanely worth it.

Check the sources below out for some good reading material on how to get started:

1. [nix.dev](https://nix.dev/)
2. [NixOS and Flakes Book](https://nixos-and-flakes.thiscute.world/)
3. [Vimjoyer's Youtube Channel](https://www.youtube.com/@vimjoyer)

In addition, I'll write another post detailing how I set it up for myself.
