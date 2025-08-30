package main

import "flag"

type Config struct {
	SrcDir    string
	OutDir    string
	PublicDir string
	TplDir    string
	Verbose   bool
}

func Load() *Config {
	src := flag.String("src", "site", "directory for markdown source files")
	out := flag.String("out", "out", "output directory for generated site")
	public := flag.String("public", "public", "public directory for public files")
	tpl := flag.String("templates", "templates", "directory for HTML templates")
	verbose := flag.Bool("verbose", false, "enable verbose logging")

	flag.Parse()

	return &Config{
		SrcDir:    *src,
		OutDir:    *out,
		PublicDir: *public,
		TplDir:    *tpl,
		Verbose:   *verbose,
	}
}
