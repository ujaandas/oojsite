package main

import (
	"fmt"
	"html/template"
)

/*
Handle everything pertaining to our final outputted HTML, including post templates.
Despite the name, a page is _not_ simply just landing pages like `index.html`, your blog post
will also eventually become a "page" - basically anything that ends up as HTML is a page.
Additionally, all pages have access to Go's templating features, allowing end-users to
access a wide breadth of global state, things like tags, page titles, options, etc...
*/

func loadPages(tmplDir, siteDir string) (*template.Template, error) {
	tmpls := template.New("")

	tmpls, err := tmpls.ParseGlob(fmt.Sprintf("%s/*.html", tmplDir))
	if err != nil {
		return nil, err
	}

	tmpls, err = tmpls.ParseGlob(fmt.Sprintf("%s/*.html", siteDir))
	if err != nil {
		return nil, err
	}

	return tmpls, nil
}
