package assets

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

func Templates() (*template.Template, error) {
	return template.ParseFS(templateFS, "templates/*.html")
}

func Static() (fs.FS, error) {
	return fs.Sub(staticFS, "static")
}
