package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gerifield/mnb-qr-go/src/server"
)

func main() {
	listen := flag.String("listen", ":8080", "HTTP listen address")
	flag.Parse()

	s := server.New()

	http.HandleFunc("/", s.GenerateHandler)

	log.Println("Listening on", *listen)
	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
