package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"
)

const PRODUCTION_BUILD = true

var tmpl *template.Template

func main() {
	var flag = *flag.Int("port", 8250, "host port")
	var listen = fmt.Sprintf("0.0.0.0:%d", flag)

	var err error
	tmpl, err = template.ParseGlob("templates/*")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Hosting at %s\n", listen)

	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	http.HandleFunc("/", handle_request)

	log.Fatal(http.ListenAndServe(listen, nil))
}

func handle_request(w http.ResponseWriter, r *http.Request) {
	if PRODUCTION_BUILD {
		var err error
		tmpl, err = template.ParseGlob("templates/*")

		if err != nil {
			log.Fatal(err)
		}
	}

	resource := "index.html"

	if r.URL.Path != "/" {
		resource = fmt.Sprintf("%s.html", path.Base(r.URL.Path))
	}

	log.Printf("Web request for %s\n", resource)

	var header bytes.Buffer
	tmpl.ExecuteTemplate(&header, "navbar.html", nil)

	var footer bytes.Buffer
	tmpl.ExecuteTemplate(&footer, "footer.html", nil)

	var head bytes.Buffer
	tmpl.ExecuteTemplate(&head, "head.html", nil)

	data := struct {
	  Navbar string
    Footer string
	  Head   string
    Error  string
	} {
	  Navbar: header.String(),
	  Footer: footer.String(),
	  Head:   head.String(),
    Error:  "None",
	}

	if err := tmpl.ExecuteTemplate(w, resource, data); err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
    data.Error = err.Error();

    tmpl.ExecuteTemplate(w, "error.html", data);

		fmt.Printf("Error: %s\n", err)
	}
}
