package model

import "html/template"

type Post struct {
	SourcePath  string
	OutputRel   string
	Filepath    string
	Frontmatter map[string]interface{}
	Snippet     string
	Content     template.HTML
	Raw         []byte
}

type GlobalData struct {
	Posts []Post
}

type TemplateData struct {
	Content     template.HTML
	Frontmatter map[string]interface{}
	Global      GlobalData
}

type PageData struct {
	Global GlobalData
}
