package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/windevkay/hardpassv2/internal/entities"
)

type templateData struct {
	CurrentYear int
	Password    *entities.Password
	Passwords   []*entities.Password
}

func formattedDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// template functions
var functions = template.FuncMap{
	"formattedDate": formattedDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// parse base template + specify template functions
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}
		// parse partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}
		// parse page
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
