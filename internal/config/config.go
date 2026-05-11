package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
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

	flag.StringVar(&cfg.OutDir, "outDir", "out", "Path to generate site in")
	flag.StringVar(&cfg.PageDir, "pageDir", "site", "Path to pages folder")
	flag.StringVar(&cfg.PostDir, "postDir", "posts", "Path to posts folder")
	flag.StringVar(&cfg.StaticDir, "staticDir", "static", "Path to static folder")
	flag.StringVar(&cfg.TemplateDir, "templateDir", "templates", "Path to templates folder")
	flag.StringVar(&cfg.ComponentDir, "componentDir", "components", "Path to components folder")
	flag.StringVar(&cfg.BaseURL, "baseUrl", "baseUrl", "Base site URL")
	flag.BoolVar(&cfg.Dev, "dev", false, "Start development server")

	flag.Parse()

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
