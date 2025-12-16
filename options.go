package main

import (
	"flag"
	"fmt"
	"os"
)

/*
All global state about where things are (ie; paths to templates, static files, pages, etc...)
are stored here, as well as any command-line flags. Later, this will be able to read from a configuration
file of sorts for better portability.
*/

type Config struct {
	outDir      string
	pageDir     string
	postDir     string
	staticDir   string
	templateDir string
}

func parseFlags() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.outDir, "outDir", "out", "Path to generate site in")
	flag.StringVar(&cfg.pageDir, "pageDir", "site", "Path to pages folder")
	flag.StringVar(&cfg.postDir, "postDir", "site/posts", "Path to posts folder")
	flag.StringVar(&cfg.staticDir, "staticDir", "static", "Path to static folder")
	flag.StringVar(&cfg.templateDir, "templateDir", "templates", "Path to templates folder")

	flag.Parse()
	// Validate all directories in one helper
	if err := validateDirs(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate that every directory exists.
func validateDirs(cfg *Config) error {
	// All required dirs
	// TODO: Can I use reflection here?
	dirs := []string{cfg.outDir, cfg.pageDir, cfg.postDir, cfg.staticDir, cfg.templateDir}

	// Ensure they all exist, create if not.
	for _, path := range dirs {
		if err := ensureDir(path); err != nil {
			return fmt.Errorf("failed to prepare directory %s: %w", path, err)
		}
	}

	return nil
}

// Check all our required directories exist. If not, create them.
func ensureDir(path string) error {
	// Check permissions
	stat, err := os.Stat(path)

	// Create if missing
	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}

	// Now check errors to avoid "no such file or directory" error if not already created
	if err != nil {
		return err
	}

	// Check it really is a directory
	if !stat.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", path)
	}
	return nil
}
