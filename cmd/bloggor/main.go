package main

import (
	"bloggor/internal/config"
	"bloggor/internal/content"
	"bloggor/internal/generate"
	"bloggor/internal/render"
	"log"
	"os"
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
