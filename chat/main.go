package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/go_blueprints/chat/client"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("../templates",
			t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "Address of the application")
	flag.Parse()
	r := client.NewRoom()
	// Handle request arriving on /
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../templates/js"))))
	http.Handle("/room", r)
	go r.Run()
	log.Printf("Starting web server on [%s]", *addr)
	// Start a webserver on 8080
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
