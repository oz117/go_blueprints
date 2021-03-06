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

	"github.com/oz117/go_blueprints/chapter_3/chat/auth"
	"github.com/oz117/go_blueprints/chapter_3/chat/client"
	"github.com/oz117/go_blueprints/chapter_3/chat/trace"
	"github.com/oz117/go_blueprints/chapter_3/chat/utils"
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
	config := utils.NewConfiguration()
	r := client.NewRoom(client.UseGravatarAvatar)
	if strings.Compare(*verbose, "true") == 0 {
		r.Tracer = trace.New(os.Stdout)
		config.Tracer = trace.New(os.Stdout)
	}
	providers, err := config.GetProviderConfig()
	if err != nil {
		log.Fatal(err)
	}
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		google.New((*providers)[0].AppID, (*providers)[0].Secret,
			(*providers)[0].Callback),
	)
	// Handle request arriving on /
	http.Handle("/", auth.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/logout", auth.Logout)
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", utils.UploaderHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../templates/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../templates/css"))))
	http.Handle("/room", r)
	http.HandleFunc("/auth/", auth.LoginHandler)
	go r.Run()
	log.Printf("Starting web server on [%s]", *addr)
	// Start a webserver on 8080
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
