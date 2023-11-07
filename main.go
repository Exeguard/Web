package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"
)

var tmpl *template.Template

func main() {
  var flag = *flag.Int("port", 8080, "host port");
  var listen = fmt.Sprintf("0.0.0.0:%d", flag);

  var err error;
  tmpl, err = template.ParseGlob("templates/*");

  if err != nil {
    log.Fatal(err);
  }

  log.Printf("Hosting at %s", listen);

  http.Handle(
    "/static/",
    http.StripPrefix(
      "/static/",
      http.FileServer(http.Dir("static")),
    ),
  );

  http.HandleFunc("/", handle_request);

  log.Fatal(http.ListenAndServe(listen, nil));
}

func handle_request(w http.ResponseWriter, r *http.Request) {
  resource := "index.html"

  if r.URL.Path != "/" {
    resource = fmt.Sprintf("%s.html", path.Base(r.URL.Path))
  }

  log.Printf("Web request for %s", resource);

  if err := tmpl.ExecuteTemplate(w, resource, nil); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    fmt.Printf("Error: %s", err)
  }
}
