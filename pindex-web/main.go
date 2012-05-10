package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	httpHost *string = flag.String("host", "0.0.0.0", "HTTP host to bind to")
	httpPort *int    = flag.Int("port", 8585, "HTTP port to bind to")
)

func main() {
	flag.Parse()

	staticDirs := []string{"js", "css", "img", "ico", "data"}
	for _, d := range staticDirs {
		route := "/" + d + "/"
		strip := "/" + d
		serve := "./" + d + "/"
		http.Handle(route, http.StripPrefix(strip, http.FileServer(http.Dir(serve))))
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	endpoint := fmt.Sprintf("%s:%d", *httpHost, *httpPort)
	log.Fatalf("%s", http.ListenAndServe(endpoint, nil))
}
