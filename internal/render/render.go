package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/prashant9154/Booking_System/internal/config"
	"github.com/prashant9154/Booking_System/internal/models"
)

var app *config.AppConfig
var pathToTemplates = "./templates"
var functions = template.FuncMap{}

// NewRenderer set the config for templates package
func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	return td
}

// render template without caching ->

// func Templates(w http.ResponseWriter, hbs string) {
// 	parsedTemplate, _ := template.ParseFiles("./templates/"+hbs, "./templates/base.layout.hbs")
// 	err := parsedTemplate.Execute(w, nil)

// 	if err != nil {
// 		fmt.Println("error parsing templates: ", err)
// 		return
// 	}
// }

// rendering template with the SIMPLE caching of templates ->

// var tc = make(map[string]*template.Template)

// func RenderTemplates(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	//check if template is already present in the cache
// 	_, inMap := tc[t]

// 	if !inMap {
// 		// need to create a template and add it to the map
// 		fmt.Println("Creating template and adding it to the cache...")
// 		err = createTemplateCache(t)

// 		if err != nil {
// 			fmt.Println("Error: ", err)
// 		}
// 	} else {

// 		// we have template in the map
// 		fmt.Println("using cached template....")
// 	}

// 	tmpl = tc[t]

// 	err = tmpl.Execute(w, nil)

// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 	}
// }

// func createTemplateCache(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.layout.hbs",
// 	}

// 	// parse the templates
// 	tmpl, err := template.ParseFiles(templates...)

// 	if err != nil {
// 		return err
// 	}

// 	// add template to cache
// 	tc[t] = tmpl

// 	return nil
// }

// rendering template with the ADVANCED caching of templates ->

func Templates(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	// create a template cache
	var tc map[string]*template.Template
	if app.UseCache {
		// get the template cache from the config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// tc, err := CreateTemplateCache()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		// log.Fatal("could not get template from template cache")
		return errors.New("could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
		return err
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.hbs", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.hbs", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.hbs", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
