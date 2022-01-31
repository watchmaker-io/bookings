package render

import (
	"bytes"
	"github.com/watchmaker-io/bookings/pkg/config"
	"github.com/watchmaker-io/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(writer http.ResponseWriter, htmlTemplate string, data *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, exists := tc[htmlTemplate]
	if !exists {
		log.Fatal("Cannot get template from template cache")
	}

	buf := new(bytes.Buffer)

	data = AddDefaultData(data)

	_ = t.Execute(buf, data)

	_, error := buf.WriteTo(writer)
	if error != nil {
		log.Println("error writing template to http response")
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, error := filepath.Glob("./templates/*.page.tmpl")
	if error != nil {
		return cache, error
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, error := template.New(name).Funcs(functions).ParseFiles(page)
		if error != nil {
			return cache, error
		}

		matches, error := filepath.Glob("./templates/*.layout.tmpl")
		if error != nil {
			return cache, error
		}

		if len(matches) > 0 {
			_, error := ts.ParseGlob("./templates/*.layout.tmpl")
			if error != nil {
				return cache, error
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
