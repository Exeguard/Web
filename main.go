package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const PRODUCTION_BUILD = false

var utilTmpl *template.Template
var tmpl *template.Template

func load_templates() {
	var files []string
	var err error

	err = filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a file and not in the "utils/" subdirectory
		if !info.IsDir() && !strings.Contains(path, "utils/") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Reading templates: %s", err)
	}

	tmpl, err = template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}

	utilTmpl, err = template.ParseGlob("templates/utils/*")

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var flag = *flag.Int("port", 8250, "host port")
	var listen = fmt.Sprintf("0.0.0.0:%d", flag)

	load_templates()

	log.Printf("Hosting at %s\n", listen)

	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)

  http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/text/robots.txt")
  })

	http.HandleFunc("/", handle_request)

	log.Fatal(http.ListenAndServe(listen, nil))
}

func handle_request(w http.ResponseWriter, r *http.Request) {
	if !PRODUCTION_BUILD {
		load_templates()
	}

	resource := "index.html"

	if r.URL.Path != "/" {
		resource = fmt.Sprintf("%s.html", path.Base(r.URL.Path))
	}

	log.Printf("Web request for %s\n", resource)

	var header bytes.Buffer
	utilTmpl.ExecuteTemplate(&header, "navbar.html", nil)

	var footer bytes.Buffer
	utilTmpl.ExecuteTemplate(&footer, "footer.html", nil)

	var head bytes.Buffer
	utilTmpl.ExecuteTemplate(&head, "head.html", nil)

	data := struct {
		Navbar string
		Footer string
		Head   string
		Error  string
	}{
		Navbar: header.String(),
		Footer: footer.String(),
		Head:   head.String(),
		Error:  "None",
	}

	if err := tmpl.ExecuteTemplate(w, resource, data); err != nil {
		data.Error = err.Error()

		tmpl.ExecuteTemplate(w, "error.html", data)

		fmt.Printf("Error: %s\n", err)
	}
}
