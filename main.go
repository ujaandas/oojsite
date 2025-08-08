package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func main() {
	tmpl := template.Must(template.ParseFiles("layouts/layout.html"))

	outputPath := filepath.Join("public", "output.html")
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	data := struct {
		Title string
	}{
		Title: "hello wurld",
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("Template saved to", outputPath)
}
