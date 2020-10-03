package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	listen := flag.String("listen", ":8080", "HTTP listen address")
	flag.Parse()

	http.HandleFunc("/", generateHandler)

	log.Println("Listening on", *listen)
	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
