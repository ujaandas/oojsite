package app

import (
	"fmt"
	"log"
	"net/http"

	"oojsite/internal/assets"
	"oojsite/internal/config"
	"oojsite/internal/content"
	"oojsite/internal/templates"
)

func Run() error {
	log.Println("Parsing options...")
	cfg, err := config.Parse()
	if err != nil {
		return err
	}
	log.Println("Options parsed!")

	log.Println("Loading templates...")
	tmpls, err := templates.Load(cfg.TemplateDir, cfg.ComponentDir, cfg.PageDir)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}
	log.Println("Templates loaded!")

	log.Println("Loading posts...")
	posts, err := content.LoadPosts(cfg.PostDir)
	if err != nil {
		return fmt.Errorf("failed to load posts: %w", err)
	}
	log.Printf("Loaded %d posts!", len(posts))

	log.Println("Rendering posts...")
	if err := content.RenderPosts(posts, cfg.OutDir, tmpls); err != nil {
		return fmt.Errorf("failed to render posts: %w", err)
	}

	log.Println("Rendering pages...")
	if err := content.RenderPages(cfg.PageDir, cfg.OutDir, posts, tmpls); err != nil {
		return fmt.Errorf("failed to render pages: %w", err)
	}

	log.Println("Building TailwindCSS...")
	if err := assets.BuildTailwind(cfg.OutDir, cfg.StaticDir); err != nil {
		return fmt.Errorf("failed to build TailwindCSS: %w", err)
	}
	log.Println("TailwindCSS built!")

	log.Println("Copying static files...")
	if err := assets.CopyStaticContents(cfg.StaticDir, fmt.Sprintf("%s/static", cfg.OutDir)); err != nil {
		return fmt.Errorf("failed to copy static files: %w", err)
	}
	log.Println("Copied static files!")

	log.Println("Building sitemap...")
	if err := assets.BuildSitemap(cfg.BaseURL, cfg.OutDir); err != nil {
		return fmt.Errorf("failed to build sitemap: %w", err)
	}
	log.Println("Built sitemap!")

	if !cfg.Dev {
		return nil
	}

	log.Println("Server started on localhost:8000!")
	return http.ListenAndServe(":8000", http.FileServer(http.Dir(cfg.OutDir)))
}
