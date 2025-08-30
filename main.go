package main

import (
	"log"
	"os"
)

func main() {
	cfg := Load()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	if cfg.Verbose {
		logger.Printf("config: %+v\n", cfg)
	}

	posts, err := DiscoverPosts(cfg.SrcDir)
	if err != nil {
		logger.Fatalf("discover posts: %v", err)
	}

	tpl, err := LoadTemplates(cfg.TplDir)
	if err != nil {
		logger.Fatalf("load templates: %v", err)
	}

	if err := GenerateSite(posts, cfg.OutDir, tpl, logger); err != nil {
		logger.Fatalf("generate site: %v", err)
	}

	logger.Println("Site generated at", cfg.OutDir)
}
