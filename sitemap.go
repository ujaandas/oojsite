package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Hold each entry
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

func generateURL(baseUrl, outDir, filePath string) string {
	rel, err := filepath.Rel(outDir, filePath)
	if err != nil {
		return baseUrl
	}

	rel = filepath.ToSlash(rel)

	// strip index.html → directory
	if before, ok := strings.CutSuffix(rel, "index.html"); ok {
		rel = before
	} else if before0, ok0 := strings.CutSuffix(rel, ".html"); ok0 {
		// strip .html → clean URL
		rel = before0
		rel += "/"
	}

	// ensure leading slash
	if !strings.HasPrefix(rel, "/") {
		rel = "/" + rel
	}

	return strings.TrimRight(baseUrl, "/") + rel
}

func addFilesToSitemap(baseUrl, outDir string, sitemap *Sitemap) error {
	err := filepath.Walk(outDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// only add html pages
		if info.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		url := generateURL(baseUrl, outDir, path)
		urlEntry := URL{
			Loc:        url,
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "monthly",
			Priority:   0.5,
		}

		sitemap.URLs = append(sitemap.URLs, urlEntry)
		return nil
	})
	return err
}

func generateSitemapFile(outDir string, sitemap *Sitemap) error {
	file, err := os.Create(fmt.Sprintf("%s/static/sitemap.xml", outDir))
	if err != nil {
		return err
	}
	defer file.Close()

	// marshal sitemap structure into XML
	enc := xml.NewEncoder(file)
	enc.Indent("", "  ")
	return enc.Encode(sitemap)
}

func buildSitemap(baseUrl, outDir string) error {
	sitemap := Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	// add files to the sitemap
	if err := addFilesToSitemap(baseUrl, outDir, &sitemap); err != nil {
		log.Fatalf("Error while adding files to sitemap: %v\n", err)
		return err
	}

	// generate the sitemap
	if err := generateSitemapFile(outDir, &sitemap); err != nil {
		log.Fatalf("Error while generating sitemap file: %v\n", err)
		return err
	}

	return nil
}
