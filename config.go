package main

import "flag"

type Config struct {
	SrcDir    string
	OutDir    string
	StaticDir string
	TplDir    string
	Verbose   bool
}

func Load() *Config {
	src := flag.String("src", "assets/content", "directory for markdown source files")
	out := flag.String("out", "public", "output directory for generated site")
	static := flag.String("static", "static", "static directory for static files")
	tpl := flag.String("templates", "templates", "directory for HTML templates")
	verbose := flag.Bool("verbose", false, "enable verbose logging")

	flag.Parse()

	return &Config{
		SrcDir:    *src,
		OutDir:    *out,
		StaticDir: *static,
		TplDir:    *tpl,
		Verbose:   *verbose,
	}
}
