package main

import (
	"encoding/json"
	"flag"
	"github.com/gobuffalo/packr"
	"io"
	"log"
	"net/http"
)

type getStatusBody struct {
	URL string `json:"url"`
}

var client http.Client

func getStatusHandler(w http.ResponseWriter, r *http.Request) {
	b := &getStatusBody{}
	if err := json.NewDecoder(r.Body).Decode(b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req, err := http.NewRequest(http.MethodGet, b.URL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = io.Copy(w, res.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	box := packr.NewBox("./static")
	port := flag.String("p", "8100", "port to serve on")
	flag.Parse()

	log.Printf("Serving %s on HTTP port: %s\n", ".", *port)
	http.Handle("/getStatus", http.HandlerFunc(getStatusHandler))
	http.Handle("/", http.FileServer(box))
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
