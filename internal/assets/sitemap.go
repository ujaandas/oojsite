package assets

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

func BuildSitemap(baseURL, outDir string) error {
	sitemap := Sitemap{XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	if err := addFilesToSitemap(baseURL, outDir, &sitemap); err != nil {
		return err
	}
	return generateSitemapFile(outDir, &sitemap)
}

func generateURL(baseURL, outDir, filePath string) string {
	rel, err := filepath.Rel(outDir, filePath)
	if err != nil {
		return baseURL
	}

	rel = filepath.ToSlash(rel)
	if before, ok := strings.CutSuffix(rel, "index.html"); ok {
		rel = before
	} else if before, ok := strings.CutSuffix(rel, ".html"); ok {
		rel = before + "/"
	}
	if !strings.HasPrefix(rel, "/") {
		rel = "/" + rel
	}
	return strings.TrimRight(baseURL, "/") + rel
}

func addFilesToSitemap(baseURL, outDir string, sitemap *Sitemap) error {
	return filepath.Walk(outDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		sitemap.URLs = append(sitemap.URLs, URL{
			Loc:        generateURL(baseURL, outDir, path),
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "monthly",
			Priority:   0.5,
		})
		return nil
	})
}

func generateSitemapFile(outDir string, sitemap *Sitemap) error {
	file, err := os.Create(fmt.Sprintf("%s/static/sitemap.xml", outDir))
	if err != nil {
		return err
	}
	defer file.Close()

	enc := xml.NewEncoder(file)
	enc.Indent("", "  ")
	return enc.Encode(sitemap)
}
