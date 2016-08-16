package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go_blueprints/chat/auth"
	"github.com/go_blueprints/chat/client"
	"github.com/go_blueprints/chat/trace"
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
	var verbose = flag.String("V", "false", "Display more information about what's going on")
	flag.Parse()
	r := client.NewRoom()
	if strings.Compare(*verbose, "true") == 0 {
		r.Tracer = trace.New(os.Stdout)
	}
	// Handle request arriving on /
	http.Handle("/", auth.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../templates/js"))))
	http.Handle("/room", r)
	go r.Run()
	log.Printf("Starting web server on [%s]", *addr)
	// Start a webserver on 8080
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
