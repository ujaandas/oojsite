package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	AllDir       string
	OutDir       string
	PageDir      string
	PostDir      string
	StaticDir    string
	TemplateDir  string
	ComponentDir string
	BaseURL      string
	Dev          bool
}

func Parse() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.AllDir, "allDir", "", "Base directory to prepend to other paths (site, posts, templates, components, static)")
	flag.StringVar(&cfg.OutDir, "outDir", "out", "Path to generate site in")
	flag.StringVar(&cfg.PageDir, "pageDir", "site", "Path to pages folder")
	flag.StringVar(&cfg.PostDir, "postDir", "posts", "Path to posts folder")
	flag.StringVar(&cfg.StaticDir, "staticDir", "static", "Path to static folder")
	flag.StringVar(&cfg.TemplateDir, "templateDir", "templates", "Path to templates folder")
	flag.StringVar(&cfg.ComponentDir, "componentDir", "components", "Path to components folder")
	flag.StringVar(&cfg.BaseURL, "baseUrl", "baseUrl", "Base site URL")
	flag.BoolVar(&cfg.Dev, "dev", false, "Start development server")

	flag.Parse()

	// Apply allDir prefix to paths that still have their default values
	if cfg.AllDir != "" {
		if cfg.PageDir == "site" {
			cfg.PageDir = filepath.Join(cfg.AllDir, "site")
		}
		if cfg.PostDir == "posts" {
			cfg.PostDir = filepath.Join(cfg.AllDir, "posts")
		}
		if cfg.StaticDir == "static" {
			cfg.StaticDir = filepath.Join(cfg.AllDir, "static")
		}
		if cfg.TemplateDir == "templates" {
			cfg.TemplateDir = filepath.Join(cfg.AllDir, "templates")
		}
		if cfg.ComponentDir == "components" {
			cfg.ComponentDir = filepath.Join(cfg.AllDir, "components")
		}
	}

	if err := validateDirs(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateDirs(cfg *Config) error {
	dirs := []string{cfg.OutDir, cfg.PageDir, cfg.PostDir, cfg.StaticDir, cfg.TemplateDir, cfg.ComponentDir}

	for _, path := range dirs {
		if err := ensureDir(path); err != nil {
			return fmt.Errorf("failed to prepare directory %s: %w", path, err)
		}
	}

	return os.RemoveAll(cfg.OutDir)
}

func ensureDir(path string) error {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", path)
	}
	return nil
}
