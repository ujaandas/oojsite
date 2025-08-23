package main

import (
	"log"
	"os"

	"oojsite/internal/config"
	"oojsite/internal/content"
	"oojsite/internal/generate"
	"oojsite/internal/render"
)

func main() {
	cfg := config.Load()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	if cfg.Verbose {
		logger.Printf("config: %+v\n", cfg)
	}

	posts, err := content.DiscoverPosts(cfg.SrcDir)
	if err != nil {
		logger.Fatalf("discover posts: %v", err)
	}

	if err := generate.CleanOutput(cfg.OutDir); err != nil {
		logger.Fatalf("clean output: %v", err)
	}

	tpl, err := render.LoadTemplates(cfg.TplDir)
	if err != nil {
		logger.Fatalf("load templates: %v", err)
	}

	if err := generate.GenerateSite(posts, cfg.OutDir, tpl, logger); err != nil {
		logger.Fatalf("generate site: %v", err)
	}

	logger.Println("Site generated at", cfg.OutDir)
}
