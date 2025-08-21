package main

import (
	"html/template"
	"log"
	"os"
)

type Page struct {
	Title string
	Body  string
}

func main() {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Parse error:", err)
	}

	if err := os.RemoveAll("public"); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll("public", 0755); err != nil {
		log.Fatal("MkdirAll error:", err)
	}

	out, err := os.Create("public/index.html")
	if err != nil {
		log.Fatal("Create file error:", err)
	}
	defer out.Close()

	page := Page{
		Title: "Hello, World!",
		Body:  "Hello, World!",
	}

	if err := tmpl.Execute(out, page); err != nil {
		log.Fatal("Execute error:", err)
	}

	log.Println("ok")
}
