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

	"github.com/go_blueprints/chapter_2/chat/auth"
	"github.com/go_blueprints/chapter_2/chat/client"
	"github.com/go_blueprints/chapter_2/chat/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
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
	data := map[string]interface{}{
		"host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "Address of the application")
	var verbose = flag.String("V", "false", "Display more information about what's going on")
	flag.Parse()
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		google.New("id",
			"your_key", "http://localhost:8080/auth/callback/google"),
	)
	r := client.NewRoom()
	if strings.Compare(*verbose, "true") == 0 {
		r.Tracer = trace.New(os.Stdout)
	}
	// Handle request arriving on /
	http.Handle("/", auth.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../templates/js"))))
	http.Handle("/room", r)
	http.HandleFunc("/auth/", auth.LoginHandler)
	go r.Run()
	log.Printf("Starting web server on [%s]", *addr)
	// Start a webserver on 8080
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
