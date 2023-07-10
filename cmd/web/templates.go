package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/windevkay/hardpassv2/internal/entities"
	"github.com/windevkay/hardpassv2/ui"
)

type templateData struct {
	CurrentYear     int
	Password        *entities.Password
	Passwords       []*entities.Password
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken	   string
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

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		
		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
