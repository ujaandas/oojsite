package config

import "flag"

type Config struct {
	SrcDir  string
	OutDir  string
	TplDir  string
	Verbose bool
}

func Load() *Config {
	src := flag.String("src", "content", "directory for markdown source files")
	out := flag.String("out", "public", "output directory for generated site")
	tpl := flag.String("templates", "templates", "directory for HTML templates")
	verbose := flag.Bool("verbose", false, "enable verbose logging")

	flag.Parse()

	return &Config{
		SrcDir:  *src,
		OutDir:  *out,
		TplDir:  *tpl,
		Verbose: *verbose,
	}
}
