package main

import (
	"log"
	"net/http"
)

func main() {
	// Handle request arriving on /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
      <html>
        <head>
          <title>Chat</title>
        </head>
        <body>
          Let's ... play a game
        </body>
      </html>
      `))
	})
	// Start a webserver on 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
