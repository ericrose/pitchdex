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
		log.Printf(
			"serving client %s (via %s)",
			r.RemoteAddr,
			func() string {
				if r.Referer() == "" {
					return "direct"
				}
				return r.Referer()
			}(),
		)
		http.ServeFile(w, r, "index.html")
	})

	endpoint := fmt.Sprintf("%s:%d", *httpHost, *httpPort)
	log.Printf("serving on %s", endpoint)
	log.Fatalf("%s", http.ListenAndServe(endpoint, nil))
}
